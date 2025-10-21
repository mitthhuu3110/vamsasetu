package com.vamsasetu.familytree;

import org.springframework.data.neo4j.repository.Neo4jRepository;

public interface RelationshipRepository extends Neo4jRepository<Relationship, String> {
}
