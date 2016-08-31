vendorutils
-----------


Utilities(or hacks pending viewpoint) to assist managing the Go `vendor/` standard while the community stabalizes around formal tooling. This isn't a long term project but if you're between [govendor](https://github.com/kardianos/govendor) and a different dependency management tool(in this case [glock](https://github.com/robfig/glock)) this library and `magazine` might be useful.

## magazine

A tool to convert `govendor` `vendor/vendor.json` file to the `GLOCKFILE` standard. Interestingly scanning the same directory `govendor` actually produces a different list of packages since it flattens dependencies which exist in external `vendor/` directories. So far this has been a good thing, however it's something to be aware of and might cause issues. Example/issue of the behavior can be found [here](https://github.com/kardianos/govendor/issues/207).

#### magazine usage

`cd magazine; go install`

Scan the `gowrapmx4j/vendor/vendor.json` file and write the GLOCKFILE which `glock` can use to sync dependencies in the GOPATH.  

`magazine -dirPath=$GOPATH/src/github.com/lytics/gowrapmx4j`

Use Glock to sync the GOPATH

`glock sync github.com/lytics/gowrapmx4j`

