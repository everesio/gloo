package flagutils

import (
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/spf13/pflag"
)

func AddRouteFlags(set *pflag.FlagSet, route *options.InputRoute) {
	set.Uint32VarP(&route.InsertIndex, "index", "x", 0, "index in the virtual service "+
		"route list where to insert this route. routes after it will be shifted back one")

	set.StringVarP(&route.Matcher.PathExact, "path-exact", "e", "", "exact path to match route")
	set.StringVarP(&route.Matcher.PathRegex, "path-regex", "r", "", "regex matcher for route. "+
		"note: only one of path-exact, path-regex, or path-prefix should be set")
	set.StringVarP(&route.Matcher.PathPrefix, "path-prefix", "p", "", "path prefix to match route")
	set.StringSliceVarP(&route.Matcher.Methods, "method", "m", []string{},
		"the HTTP methods (GET, POST, etc.) to match on the request. if empty, all methods will match ")
	set.StringSliceVarP(&route.Matcher.HeaderMatcher.Entries, "header", "d", []string{},
		"headers to match on the request. values can be specified using regex strings")

	set.StringVarP(&route.Destination.Upstream.Name, "dest-name", "u", "",
		"name of the destination upstream for this route")
	set.StringVarP(&route.Destination.Upstream.Namespace, "dest-namespace", "s", defaults.GlooSystem,
		"namespace of the destination upstream for this route")

	set.StringVarP(&route.Destination.DestinationSpec.Aws.LogicalName, "aws-function-name", "a", "",
		"logical name of the AWS lambda to invoke with this route. use if destination is an AWS upstream")
	set.BoolVarP(&route.Destination.DestinationSpec.Aws.ResponseTransformation, "aws-unescape", "", false,
		"unescape JSON returned by this lambda function (useful if the response is not intended to be JSON formatted, "+
			"e.g. in the case of static content (images, HTML, etc.) being served by Lambda")

	set.StringVarP(&route.Destination.DestinationSpec.Rest.FunctionName, "rest-function-name", "f", "",
		"name of the REST function to invoke with this route. use if destination has a REST service spec")
	set.StringSliceVar(&route.Destination.DestinationSpec.Rest.Parameters.Entries, "rest-parameters", nil,
		"Parameters for the rest function that are to be read off of incoming request headers. format specified as follows: "+
			"'header_name=extractor_string' where header_name is the HTTP2 equivalent header (':path' for HTTP 1 path).\n\n"+
			"For example, to extract the variable 'id' from the following request path /users/1, where 1 is the id:\n"+
			"--rest-parameters ':path='/users/{id}'")

	set.Var(&route.Plugins.PrefixRewrite, "prefix-rewrite", "rewrite the matched portion of HTTP requests with this prefix.\n"+
		"note that this will be overridden if your routes point to function destinations")
}

func RemoveRouteFlags(set *pflag.FlagSet, route *options.RemoveRoute) {
	set.Uint32VarP(&route.RemoveIndex, "index", "x", 0, "remove the route with this index in the virtual service "+
		"route list")
}
