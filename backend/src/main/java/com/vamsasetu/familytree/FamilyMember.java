package com.vamsasetu.familytree;

import org.springframework.data.neo4j.core.schema.Id;
import org.springframework.data.neo4j.core.schema.Node;
import java.time.LocalDate;

@Node("FamilyMember")
public class FamilyMember {
    @Id
    private String id;
    private String name;
    private String gender;
    private LocalDate dateOfBirth;
    private String metaData;
    // Add more attributes for full Indian family details

    // getters and setters
}
