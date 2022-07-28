// *** WARNING: this file was generated by test. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Myedgeorder.Inputs
{

    /// <summary>
    /// Configuration filters
    /// </summary>
    public sealed class ConfigurationFiltersArgs : global::Pulumi.ResourceArgs
    {
        [Input("filterableProperty")]
        private InputList<Inputs.FilterablePropertyArgs>? _filterableProperty;

        /// <summary>
        /// Filters specific to product
        /// </summary>
        public InputList<Inputs.FilterablePropertyArgs> FilterableProperty
        {
            get => _filterableProperty ?? (_filterableProperty = new InputList<Inputs.FilterablePropertyArgs>());
            set => _filterableProperty = value;
        }

        /// <summary>
        /// Product hierarchy information
        /// </summary>
        [Input("hierarchyInformation", required: true)]
        public Input<Inputs.HierarchyInformationArgs> HierarchyInformation { get; set; } = null!;

        public ConfigurationFiltersArgs()
        {
        }
        public static new ConfigurationFiltersArgs Empty => new ConfigurationFiltersArgs();
    }
}
