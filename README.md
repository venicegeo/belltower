# BellTower

_Campanile (cam-pin-EE-lay): a bell tower or watchtower; the tallest building in Venice_


# Technical Goals

_"Do Nothing 'til You Hear from Me"_

1. Needs to support existing applications/services as they are, as no one
will update their service to meet some new system we invent (e.g. try to
get a 3rd party to do REST calls to us or monitor a message queue). In
particular, we must not require the client services to be able to maintain
state.

2. The system should not preclude any ability the apps may have to scale
themselves. The system is not responsible for scaling clients.

3. The communications protocol between workflow nodes must be simple, to make
it easier to write components. It also must be language-independent, to allow
for a future where components can be written in any language.  (In addition,
components should have no shared memory between them, so that someday we
have the option of moving to a fully distributed system.) This means all
communications should be done via JSON. 

4. We do not want to reinvent the idea of a workflow or build something from
scratch. If there is a open source workflow package that meets our criteria
and allows us to readily build on top of it, we should use it. Such an
underlying system should make no assumptions about a security model; our
system will layer that on top. If we use an open source framework, it should
be well-supported by the community _or_ it should be relatively small enough
that we could fork and own it ourselves.

5. All "local" processing (the components) should not be heavy-weight or
long-running: they should not do much more than simple JSON processing. Heavy
work should be done at the client level, invoked remotely by the local components.

6. Ideally, for development purposes, the system should be able to be run on a
laptop.

7. Our focus is on _orchestration,_ meaning taking charge over which thing
runs when and with what (lightweight) inputs. In contrast, our focus is _not_
on _data processing,_ meaning the actual computational work on the (often
heavy) data files. Furthermore, our focus is not on supporting high-speed,
high-volume, streaming data feeds, e.g. Twitter, as other, existing frameworks
do that well already.

8. The system is not responsible for "starting" or "owning" client jobs.
The system will invoke remote services, but any "job management" within
the context of that job belomgs to the client service.

9. The system should provide a rich library of components to perform the
basic operations common workflows will require, such as logging, timers,
simple JSON transformations, conditionals, etc.



# Example Use Cases

* **If This Then That (IFTTT):** When a file is updated on GitHub, send me
a note via Slack.

* **Monitoring:** Watch a set of detection feeds, S3 buckets, and web pages.
Every night, summarize the changes/updates into a single email to me. Allow
me to add/remove feeds while the system is running.

* **Data Processing/Analytics:** Every hour, check Planet feed for new imagery
and if any image is in the given AOIs, then run image through Beachfront to
compute the coast line, then report the results to me. If any error occurs,
log it and email me a nightly report of all errors. Allow me to change the
AOIs while the system is running.


# Glossary

* **Component:** a node in the _graph_ performs some function, such as writing
to a file, performing a simple computation, or invoking a remote service.
Components contain one or more input _ports_ and one or more output ports.

* **Connection:** the link formed between an output _port_ of one _component_
and the _input_ port of another component. _Messages_ are sent across connections.

* **Graph:** the set of nodes _(components)_ and edges _(connections)_ that
together make up a single _workflow._

* **Message:** the data payload passed between _components_ via their
_ports._ Messages are formatted as JSON objects.

* **Port:** an input (or output) channel to a _component_, through which a
message can be received from (or sent to) another component.

* **Workflow:** an ordered set of operations that together perform some end goal
task, such as providing notification when a file is added to an S3 bucket or
running a series of analytical operations on a piece of data. A workflow is
represented formally by a _graph._


# TO DO

## Now

* error propagation, both within the libraries and from running components
* add "verbose/logging" mode
* revive file watcher lib (and tests)
* nice model for error handling in general

## Next

* can the graph contain cycles?
* add support for pre/post conditions
* "slow motion" mode
* support running more than one graph at a time
* all components should have these fields at the core level:
  * date started
  * number of messages received, processed, etc
  * cpu and wall time used
* design the server-level system, including:
  * security
  * user management
  * database to persist state
* add DSL support - parser, etc.
* add argument type checking
* add metadata support
* add "notes" support
* add automated description generation
* add funcs to govaluate library
* implement rest of components library
* document classes
* code coverage
* linting
* "names" should only be alphanumeric
* build more infrastructure to make defining Components easier
* need a /dev/null (Grounder) component?
* validate graph connectivity
* remove START/STOP req'ments; tie all open output ports to STOP (ground?)
* design an AND component
* put panic-checks around all goflow calls (and one big one at app level?)
* can we generalize Replicator to have num output ports set at config time, e.g. an array of Output chans?
* Document the syntax extensions in govaluate
* add test cases for wrong number of port connections; goflow seems to have bad diagnostics for this,
  need to handle ourselves in validation
* Should we allow the graph to be changed while it is running? (does goflow support this?)
  

## Future

* Use CWL (http://www.commonwl.org/) to describe ommand-line usage for proxy nodes?
* Sheller components (sh, ssh) are a security hole
* Provenance tracking - collect history of processing from each component
* Components eventually become almost lamdba-like, or maybe get fully disributed using message queues, or...
* 


# Designing Components and Networks

* A component must have at least one input port. It must be read by the
  component, even if the input value is not going to be used.
* A component must have at least one output port. It must be written to
  by the component, even if the value is just the "empty payload" of "{}".

# Designing Networks

* In the network graph, each input and output port must be connected to
  another component. You may use the special START component as a dummy input
  and the special STOP component as a dummy output. _(Note: some day we will
  "ground" any unused ports automatically.)_
* An output port must be connected to one and only one input port.
* More than one output port may connect to the same input port. (This is equivalent to
  what would be an Or component, if we provided one.)


# Library of Components

(components in _italics_ are not yet implemented)

* A
  * **Adder** - just adds a set value to a field: "in.x + config.y -> out.z"
* B
  * _Beachfronter - runs a BF command_
* C
  * **Copier** - duplicates the input, e.g. "in -> out1, out2
* F
  * _FileWatcher - watches a (local) directory for changes (new files, deletes, modifies)_
* H
  * _HttpVerber - executes a HTTP GET, POST, etc, using in as the body and sending the repsonse to out_
* J
  * _JQer - runs a JQ command on the input_
  * _Joiner - waits until it has an message on all input ports, then sends the concatenation(?) to out_
* L
  * **Logger** - writes in to a file (or stdout (or stderr))
* M
  * _Mailer - sends mail, with body (and To/Subject?) taken from in_
* O
  * _Orer - when an input is recieved at any of two or more input ports, forwards the result to out_
* P
  * _Piazzaer - runs a Piazza command_
* R
  * _RabbitMQPoster - sends in to an MQ queue_
  * _RabbitMQWatcher - watches an MQ queue for new data, and sends it to out_
  * _RandomGenerator - sends random numbers, strings, etc to out_
  * **Remapper** - remaps ("changes the name of") a field: in.x -> out.y
* S
  * _S3Watcher - watches an S3 bucket for changes_
  * _Sampler - forwards every Nth input message to out (or maybe sends the "average" of the N messages)_
  * _SHer - runs a shell command, with in->stdin and stdout->out_
  * _SimpleFunction - allows for simple math, string manipulation, etc, to be done on a given field_
  * _Sleeper - sleeps for a period of time, then forwards in to out_
  * _SSHer - runs am ssh shell command, with in->stdin and stdout->out_
* T
  * **Ticker** - sends a simple output every N seconds
  * _Timer - sends a simple output at a specified time (e.g. every day at midnight)_
* W
  * _WebPageWatcher - watches a web page (web site?), sends changes to out_



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


# Credits

Belltower happily uses GoFlow (https://github.com/trustmaster/goflow) and GoValuate (https://github.com/Knetic/govaluate).
