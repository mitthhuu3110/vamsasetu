package com.vamsasetu.familytree;

import org.springframework.data.neo4j.repository.Neo4jRepository;

public interface FamilyMemberRepository extends Neo4jRepository<FamilyMember, Long> {
    // add custom queries if needed
}
