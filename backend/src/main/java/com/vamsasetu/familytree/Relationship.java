package com.vamsasetu.familytree;

import org.springframework.data.neo4j.core.schema.Id;
import org.springframework.data.neo4j.core.schema.RelationshipProperties;
import org.springframework.data.neo4j.core.schema.TargetNode;
import org.springframework.data.neo4j.core.schema.GeneratedValue;

@RelationshipProperties
public class Relationship {
    @Id
    @GeneratedValue
    private Long id;
    private String type; // e.g. father, cousin, in-law
    private FamilyMember from;
    @TargetNode
    private FamilyMember to;
    private String metaData;
    // getters and setters
}
