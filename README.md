#graveldb
persistent (on disk) graph database library for embedded use

## intro
just dabbling around with creating an embedded persistent graph database in go.

This is an exercise in learning go for me. I intend to use it as a generic backing store for other applications like note taking etc..

uses leveldb as a backing store

it sort of works right now.
it's a very basic property graph.. 
there are no optimizations for queries yet 
there is no write collision prevention, yet, so do writes serially (addvertex, addedge, deletevertex, deleteedge) - now implemented writelocks using mutex for now

just follow the graph_test.go - open a graph, add/delete vertices, add/delete edges, and query using edge following, e.g vertex.OutEdges() /vertex.InEdges() and edge.OutVertex(), edge.InVertex().. 

you can also manually check for interesting edges using edge.labels
or interesting vertices and edges using properties.



## concept
the structure (and interfaces) of the graph database tries to follow most of the blueprint graph api..


* Graph: An object that contains vertices and edges.
  * Element: An object that can have any number of key/value pairs associated with it (i.e. properties)
    * Vertex: An object that has incoming and outgoing edges.
    * Edge: An object that has a tail and head vertex.


A property graph has these elements:

1. a set of vertices
  1. each vertex has a unique identifier.
  1. each vertex has a set of outgoing edges.
  1. each vertex has a set of incoming edges.
  1. each vertex has a collection of properties defined by a map from key to value.
1. a set of edges
  1. each edge has a unique identifier.
  1. each edge has an outgoing tail vertex.
  1. each edge has an incoming head vertex.
  1. each edge has a label that denotes the type of relationship between its two vertices.
  1. each edge has a collection of properties defined by a map from key to value.


## implementation

the graph will be persisted in databases / tables on disk. These are the database descriptions

the edges are indexed in a hexastore derived index implementation


1. element 
  1. key=element id
  1. value= type (vertex or edge)

1. hexastore
  1. key = one of
    * spo::A::C::B
    * sop::A::B::C
    * ops::B::C::A
    * osp::B::A::C
    * pso::C::A::B
    * pos::C::B::A
  where 
    * A is element id for vertex as originating vertex or subject
    * B is element id for vertex as terminating vertex or object
    * C is element id for edge connecting the vertii or predicate
  1. value = label 

1. edges - mostly for quick rehydration -- will be removed later
  1. key = edge (element id)
  1. value = subject (outvertex), object (invertex), label
1. property 
  1. key = elemenid::property
  1. value = value

1. meta
  1. metadata about store
  1. some of the keys that might be used are
    1. nextid
    1. number of vertexes
    1. number of edges

  
