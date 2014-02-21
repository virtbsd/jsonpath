/*
(BSD 2-clause license)

Copyright (c) 2014, Shawn Webb
All rights reserved.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

   * Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package main

import (
    "fmt"
    "os"
    "io/ioutil"
    "github.com/bitly/go-simplejson"
    "strings"
    "strconv"
)

func resolvePathObj(j *simplejson.Json, path string) *simplejson.Json {
    if j == nil {
        return nil
    }

    if idx := strings.Index(path, "["); idx > 0 {
        var idx2 int
        var objidx int
        var err error

        idx2 = strings.Index(path, "]")

        if objidx, err = strconv.Atoi(path[idx+1:strings.Index(path, "]")]); err != nil {
            return nil
        }

        oldpath := path[0:idx]
        newpath := path[idx2+2:]

        return resolvePathObj(j.GetPath(oldpath).GetIndex(objidx), newpath)
    } else {
        if resolved := j.GetPath(path); resolved != nil {
            return resolved
        }
    }

    return nil
}

func resolvePath(j *simplejson.Json, path string) string {
    var newpath string

    if strings.HasPrefix(path, "len:") {
        newpath = path[4:]
    } else {
        newpath = path
    }

    if obj := resolvePathObj(j, newpath); obj != nil {
        if strings.HasPrefix(path, "len:") {
            if arr, err := obj.Array(); err == nil {
                return strconv.Itoa(len(arr))
            } else {
                fmt.Fprintf(os.Stderr, "[%s] ERROR: %s\n", newpath, err.Error())
            }
        } else {
            if s, err := obj.String(); err == nil {
                return s
            } else {
                fmt.Fprintf(os.Stderr, "[%s] ERROR: %s\n", path, err.Error())
            }
        }
    }

    return ""
}

func main() {
    bytes, err := ioutil.ReadAll(os.Stdin)
    if err != nil {
        panic(err)
    }

    obj, err := simplejson.NewJson(bytes)
    if err != nil {
        panic(err)
    }

    for i := 1; i < len(os.Args); i++ {
        if resolved := resolvePath(obj, os.Args[i]); len(resolved) > 0 {
            fmt.Printf("%s\n", resolved)
        }
    }
}
