// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import * as grpc from "@grpc/grpc-js";
import { log, Resource } from "..";

import * as settings from "../runtime/settings";
import { debuggablePromise } from "./debuggable";
import { deserializeProperties, deserializeProperty, serializeProperties, serializeProperty } from "./rpc";

const gstruct = require("google-protobuf/google/protobuf/struct_pb.js");
const structproto = require("google-protobuf/google/protobuf/struct_pb.js");

const callbackproto = require("../proto/callback_pb.js");
const callbackrpc = require("../proto/callback_grpc_pb.js");
const resproto = require("../proto/resource_pb.js");

let server: CallbackServer | undefined;

/** @internal */
export function shutdownServer() {
    server?.forceShutdown();
}

export async function registerCallback(callback: (args?: any) => any): Promise<string> {
    if (!server) {
        const grpcServer = new grpc.Server({
            "grpc.max_receive_message_length": settings.maxRPCMessageSize,
        });
        server = new CallbackServer(grpcServer);

        grpcServer.addService(callbackrpc.CallbackService, server);
        const port = await new Promise<number>((resolve, reject) => {
            grpcServer.bindAsync("127.0.0.1:0", grpc.ServerCredentials.createInsecure(), (err, p) => {
                if (err) {
                    reject(err);
                } else {
                    resolve(p);
                }
            });
        });
        server.setPort(port);
        grpcServer.start();
    }

    return server.registerCallback(callback);
}

class CallbackServer implements grpc.UntypedServiceImplementation {
    private port: number | undefined = undefined;
    private readonly callbacks: ((args?: any) => any)[] = [];

    constructor(private readonly grpcServer: grpc.Server) {
    }

    // Satisfy the grpc.UntypedServiceImplementation interface.
    [name: string]: any;

    public forceShutdown() {
        this.grpcServer.forceShutdown();
    }

    public setPort(port: number) {
        this.port = port;
    }

    private getCallback(reference: string): (args?: any) => any {
        const index = this.getIndexFromReference(reference);
        return this.callbacks[index];
    }

    private getIndexFromReference(reference: string): number {
        const i = reference.indexOf("/");
        if (i === -1) {
            throw new Error(`Could not determine index from callback reference "${reference}"`);
        }
        const index = reference.substring(i + 1);
        if (index.length === 0) {
            throw new Error(`No token in callback reference "${reference}"`);
        }
        return parseInt(index, 10);
    }

    public registerCallback(callback: (args?: any) => any): string {
        // Has this callback already been added?
        let index = this.callbacks.indexOf(callback);

        // If not, add it.
        if (index === -1) {
            index = this.callbacks.length;
            this.callbacks.push(callback);
        }

        // Return the callback reference <address>/<token>, where <token>
        // is an opaque string used by the SDK to lookup the function.
        // In this case, it's the index into the callbacks array.
        //
        // TODO: Consider including some extra information at the end of the token
        // like the name of the function and source location, to help with
        // debuggability.
        const reference = `${this.port}/${index}`;
        return reference;
    }

    public async invoke(call: any, callback: any): Promise<void> {
        try {
            const req: any = call.request;
            const reference: string = req.getReference();
            const args: any = req.getArgs().toJavaScript();

            const deserializedArgs = deserializeProperty(args);

            const callbackToInvoke = this.getCallback(reference);
            const result = callbackToInvoke(deserializedArgs);

            const deps = new Set<Resource>();
            const serializedResult = await serializeProperty("", result, deps, {
                keepOutputValues: true,
            });

            const resp = new callbackproto.CallbackInvokeResponse();
            resp.setReturn(structproto.Struct.fromJavaScript(serializedResult));

            callback(undefined, resp);
        } catch (e) {
            callback(e, undefined);
        }
    }
}

export function createRemoteCallback(reference: string): (args?: any) => Promise<any> {
    return async (args?: any): Promise<any> => {
        const serializedArgs = await serializeProperties("", args, {
            keepOutputValues: true,
        });

        const req = new resproto.ResourceInvokeRequest();
        req.setTok("pulumi:pulumi:invokeCallback");
        req.setArgs(gstruct.Struct.fromJavaScript({
            reference,
            args: serializedArgs,
        }));
        req.setProvider("");
        req.setVersion("");
        req.setAcceptresources(true);

        const monitor: any = settings.getMonitor();
        const resp = await debuggablePromise(new Promise<any>((resolve, reject) =>
            monitor.invoke(req, (e: grpc.ServiceError, innerResponse: any) => {
                if (e) {
                    reject(e);
                } else {
                    resolve(innerResponse);
                }
            })), "");
        return deserializeProperties(resp.getReturn());
    };
}
