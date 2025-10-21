package com.vamsasetu.familytree;

import org.springframework.data.neo4j.core.schema.Id;
import org.springframework.data.neo4j.core.schema.RelationshipProperties;
import org.springframework.data.neo4j.core.schema.TargetNode;
import org.springframework.data.neo4j.core.schema.StartNode;

@RelationshipProperties
public class Relationship {
    @Id
    private String id;
    private String type; // e.g. father, cousin, in-law
    @StartNode
    private FamilyMember from;
    @TargetNode
    private FamilyMember to;
    private String metaData;
    // getters and setters
}
