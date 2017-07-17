_Campanile (cam-pin-EE-lay): a bell tower or watchtower; the tallest building in Venice_

# Goals

1. Needs to support existing applications/services as they are, as no one will update their service to meet some new system we invent (e.g. try to get a 3rd party to do REST calls to us or monitor a message queue).

2. The system should not preclude any ability the apps may have to scale themselves. The system is not responsible for scaling clients.

3. The communications protocol between workflow nodes must be simple, to make it easier to write components. It also must be language-independent, to allow for a future where components can be written in any language. 
(In addition, components should have no shared memory between them, so they can be deployed to a distributed system.) This means all communications should be doen via JSON. 

4. We must not require the client services to be able to maintain state.

5. We do not want to reinvent the idea of a workflow or build something from scratch. If there is a open source workflow package that meets our criteria and allows us to readily build on top of it, we should use it. Such an underlying system should make no assumptions about a security model; our system will layer that on top. If we use an open source framework, it should be well-supported by the community _or_ it should be relatively small enough that we could fork and own it ourselves.

6. All "local" processing (the components) should not be heavy-weight or long-running: they should not do much more than simple JSON processing. Heavy work should be done at the client level, invoked remotely by the local components.

7. Ideally, for development purposes, the system should be able to be run on a laptop.

8. Our focus is on _orchestration,_ meaning taking charge over which thing runs when and with what (lightweight) inputs. Our focus is not on _data processing,_ meaning the actual computational work on the (often heavy) data files. Furthermore, our focus is not on supporting high-speed, high-volume, streaming data feeds, e.g. Twitter, as other systems do that well already.

9. The system is not responsible for "starting" or "owning" client jobs. The system will invoke remote services, but any "job management" within the context of that job belomgs to the client service.

10. 

Open: Should we allow the graph to be changed while it is running? (does goflow support this?)



# Example Use Cases

* **Data processing chain:** Every hour, check Planet feed for new imagery and if any image in the given AOIs, then run image through BF to get coast line, then report it to me. If any error, log it and email me a nightly report of all errors. Allow me to change the AOIs.

* **Monitoring:** Watch a set of detection feeds, S3 buckets, and web pages. Every night, summarize the changes/updates into a single email to me. Allow me to add/remove feeds.

* **If This Then That (IFTTT)** When a file is updated on GitHub, send me a note via Slack.


# TO DO

## Now

* conventions and best practices for writing components
* error propagation
* glossary of terms
* add support for pre/post conditions
* generalize executing a graph


## Next

* "slow motion" mode
* support running more than one graph at a time
* all components should have these fields at the core level:
  * date started
  * number of messages received, processed, etc
  * cpu and wall time used
  * ...
* design the server-level system, including:
  * security
  * user management
  * database to persist state

## Future

* Use CWL (http://www.commonwl.org/) to describe ommand-line usage for proxy nodes?
* Document the syntax extensions in govaluate
* Sheller components (sh, ssh) are a security hole
* Provenance tracking - collect history of processing from each component
* Components eventually become almost lamdba-like, or maybe get fully disributed using message queues, or...
* 


# Library of Components

* Ticker
* Adder
* Remapper
* ...
* FieldOperator
* Printer
* JQer
* SHer
* SSHer
* FileWatcher
* S3Watcher
* HttpVerber
* Mailer
* Logger
* RabbitMQer
* Piazzaer
* Beachfronter
* WebPageWatcher
* Joiner
* Sleeper
* RandomGenerator
* Selector
* Replicator
* Sampler


# Concept for a Graph DSL

```
# this is my graph
graph

    meta
        contact:
        version:
        description:
    endmeta
  
    components

        component
            note: later this can be changed to a Frobber2 component
            type: Frobber
            name: myfrobber
            precondition: true  # for now
            # because y can't be bigger than x
            postcondition: x >= y
            config
                x: 5
                y: "foo"
                z: struct
                    alpha: int
                    beta: int
                end struct
            endconfig
        endcomponent

    endcomponents

    connections
        connection
            from: component.port
            to: component.port
        endconnection
    endconnections
endgraph
```

Notes:
* see also, of course, the DSL for FBP (https://github.com/flowbased/fbp), which has a nice syntax for wrting the connection but doesn't allow for a decl block of the components, metadata, config blocks for components, etc
* comments start with `#` and go to end of line
* a `note` field stores a text string in the current object;
  you can have more than one note per object


