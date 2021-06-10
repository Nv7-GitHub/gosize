package main

import (
	"regexp"
	"strings"
)

const tmplpart = `<
 (?:
    [^<>]
  | (?:
     <
      (?:
        [^<>]
      | (?:
          <
            (?:
              [^<>]
            | (?:
                <
                  (?:
                    [^<>]
                  | (?:
                      <
                        (?:
                          [^<>]
                        | (?:
                            <
                              (?:
                                [^<>]
                              | (?:
                                  <
                                    (?:
                                      [^<>]
                                    | (?:
                                        <
                                          (?:
                                            [^<>]
                                          | (?:
                                              <
                                                [^<>]*
                                              >
                                            )
                                          )*
                                        >
                                      )
                                    )*
                                  >
                                )
                              )*
                            >
                          )
                        )*
                      >
                    )
                  )*
                >
              )
            )*
          >
        )
      )*
     >
    )
 )*
>`

var parengroup = strings.ReplaceAll(strings.ReplaceAll(tmplpart, "<", `\(`), ">", `\)`)

var undefre = regexp.MustCompile(`^\s*0\s+U\s+`)
var entriesre = regexp.MustCompile(`^\s*([0-9a-fA-F]+)\s+([0-9]+)\s+(\S+)\s+(.*)$`)
var cpppath = `
(?:
  \([^)]*\)
| \{[^}]*\}
| ~?
  (?:
     \$?\w+
   | operator(?:[^\(]+|\(\))
  )
  (?:` + tmplpart + `)?
  (?:` + parengroup + `)?
  (?:\sconst)?
)::
| [a-zA-Z]+_
`
var cpppathre = regexp.MustCompile(cpppath)
var cppsymre = regexp.MustCompile(`
^
# prefix
(?:guard\svariable\sfor\s)?
(

  (?:
    (?:
       \w
     | ::
     | -
     | \*
     | \&
     | (?:` + tmplpart + `)
    )+
    \s
  )*
)
(

  (?:` + cpppath + `)*
)
(
  \{[^}]*\}
|

  ~?
  (?:
     \$?\w+
   | operator(?:[^\(]+|\(\))
   | \._\d+
  )
  (?:\[[^][]*\])?
  (?:` + tmplpart + `)?
  (?:` + parengroup + `
      (?:\sconst)?
      (?:\s\[.*\])?
      (?:\s\(
          (?:\.(?:constprop|part|isra).\d+)+
      \))?
  )?
  (?:\*+)?
)
$
`)

const gopathparts = `
    \(
      (?:
        [^()]
      | \( [^()]* \)
      )*
    \)\.                  
    | struct\s\{
    (?:
        [^{}]
      | \{[^{}]*\}
    )*
    \}\.                 
    | \$?(?:\w|-|%)+\.      
    | glob\.\.             
    | \.gobytes\.    
    | (?:\w|\.|-|%)+/
`
const golastpart = `
    (?:
      (?:
          \.?
          (?:
            (?:\w|-|%)+ 
          | \( [^()]* \)
          )
          (?:-fm)?)
    | struct\s\{
        (?:
            [^{}]
          | \{[^{}]*\}
        )*
        \}
    )
    (?:,
      (?:
        (?:(?:\w|\.|-|%)+/)*
        (?:\w|-|%|\.)+
        | interface\s\{ (?: [^{}] | \{ [^{}]* \} )* \}
      )
    )?
    | initdone\.
    | initdone·
`

var gosymre = regexp.MustCompile(`
^
  (
    (?:go\.itab\.\*?)?
  )
  (
    (?:` + gopathparts + `)*
  )
  (
    ` + golastpart + `
  )
$
`)

var gopathpartsre = regexp.MustCompile(gopathparts)
