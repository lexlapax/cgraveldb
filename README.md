just dabbling around with creating an embedded persistent graph database in go.

this is an exercise in learning go for me.. 

uses leveldb as a backing store

nothing works right now..


the structure (and interfaces) of the graph database tries to follow most of the blueprint graph api..


Graph: An object that contains vertices and edges.
 Element: An object that can have any number of key/value pairs associated with it (i.e. properties)
  Vertex: An object that has incoming and outgoing edges.
  Edge: An object that has a tail and head vertex.

A property graph has these elements:

a set of vertices
 each vertex has a unique identifier.
 each vertex has a set of outgoing edges.
 each vertex has a set of incoming edges.
 each vertex has a collection of properties defined by a map from key to value.
a set of edges
 each edge has a unique identifier.
 each edge has an outgoing tail vertex.
 each edge has an incoming head vertex.
 each edge has a label that denotes the type of relationship between its two vertices.
 each edge has a collection of properties defined by a map from key to value.

databases:
element 
 key=element id
 value= type (vertex or edge)

hexastore
 key = one of
  spo::A::C::B
  sop::A::B::C
  ops::B::C::A
  osp::B::A::C
  pso::C::A::B
  pos::C::B::A
 where 
 A is element id for vertex as originating vertex or subject
 B is element id for vertex as terminating vertex or object
 C is element id for edge connecting the vertii or predicate
value = label 

property 
key = elemenid::property
value = value

meta
 metadata about store
 some of the keys that might be used are
  nextid
  number of vertexes
  number of edges

  
