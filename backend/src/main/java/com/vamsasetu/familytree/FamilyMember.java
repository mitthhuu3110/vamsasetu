package com.vamsasetu.familytree;

import org.springframework.data.neo4j.core.schema.Id;
import org.springframework.data.neo4j.core.schema.Node;
import org.springframework.data.neo4j.core.schema.GeneratedValue;
import java.time.LocalDate;

@Node("FamilyMember")
public class FamilyMember {
    @Id
    @GeneratedValue
    private Long id;
    private String name;
    private String gender;
    private LocalDate dateOfBirth;
    private String metaData;
    // Add more attributes for full Indian family details

    // getters and setters
}
