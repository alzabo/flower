# Flows

## Table of Contents

* [flowWithLotsOfComments](#flowwithlotsofcomments)
* [do-something](#do-something)
* [idk](#idk)
* [default](#default)
* [inlineCommentedFlow](#inlinecommentedflow)

## flowWithLotsOfComments

Flow that doesn't do anything

This flow will do nothing. It accepts no inputs and
has no outputs
Some sequencess may `${render}` `${strangely}`
`${more.expr}`

```yaml
- step:
    - in: this
      out: that

```

Defined in [test/f1.yaml](test/f1.yaml#L10)

[⬆️ Return to Contents](#table-of-contents) 

------


## do-something


Flow with more banners
This flow has 2 types of banners


This is the second bit
This had banners too
- Arg1
- Arg2 


```yaml
- call: idk

```

Defined in [test/f2.yaml](test/f2.yaml#L12)

[⬆️ Return to Contents](#table-of-contents) 

------


## idk


Flow with banners         

comments about the flow   
This flow doesn't do much 


```yaml
- call: default
- fail:
    - idk why

```

Defined in [test/f2.yaml](test/f2.yaml#L20)

[⬆️ Return to Contents](#table-of-contents) 

------


## default

-- default flow doc

```yaml
- task: http
  in:
    url: "https://google.com"
  out: result
- if: ${not result.ok}
  then:
    - log: "task failed: ${result.error}"

```

Defined in [test/g/g1.yml](test/g/g1.yml#L4)

[⬆️ Return to Contents](#table-of-contents) 

------


## inlineCommentedFlow

it's a comment

```yaml
- fail:

```

Defined in [test/g/g1.yml](test/g/g1.yml#L14)

[⬆️ Return to Contents](#table-of-contents) 

------


