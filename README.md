SlashedQueryUrls
================
##### This is a stable-version but not well-tested. 

A simple function to generate paths compatible with gorilla-mux(http://www.gorillatoolkit.org/pkg/mux) from a simple pattern.

### The Package

There are just one exported function

###### func GetUrlVarsPermutations(baseUrl string, isLong bool) []string 
* If isLong is false, the baseUrl must attend a pattern:
  * ^(/[^{}/ ]+)*(/{[^{}/ ]+}[?]?)*$
* If isLong is false, the baseUrl must attend a pattern:
  * ^(/[^{}/ ]+)*((/[^{}/ ]+)/{[^{}/ ]+}[?]?)*$

non-long example:

    paths := muxUrlGen.GetUrlVarsPermutations("/img/{type}/{length}?/{width}?", false)
    fmt.Println(paths) 
    /* Output:
    /img/type/{type}
    /img/type/{type}/width/{width}
    /img/type/{type}/height/{height}
    /img/width/{width}/type/{type}
    /img/height/{height}/type/{type}
    /img/type/{type}/width/{width}/height/{height}
    /img/type/{type}/height/{height}/width/{width}
    /img/width/{width}/height/{height}/type/{type}
    /img/width/{width}/type/{type/height/{height}
    /img/height/{height}/width/{width}/type/{type}
    /img/height/{height}/type/{type}/width/{width}
    */

  




##### The motivation

You could want route several paths to the same handler

* /address/street/{street}/number/{number:[0-9]+}/city/{city}
* /address/street/{street}/city/{city}/number/{number:[0-9]+}
* /address/number/{number:[0-9]+}/street/city/{city}/{street}
* /address/number/{number:[0-9]+}/street/{street}/city/{city}
* /address/city/{city}/street/{street}/number/{number:[0-9]+}
* /address/city/{city}/number/{number:[0-9]+}/street/{street}
* /address/city/{city}/street/{street}
* /address/street/{street}/city/{city}

There aren't a concise way to do that with gorilla-mux. 
I made that package to generate permutations on a url

