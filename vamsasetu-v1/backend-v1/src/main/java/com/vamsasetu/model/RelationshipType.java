package com.vamsasetu.model;

public enum RelationshipType {
    // Immediate Family
    PARENT,
    CHILD,
    SPOUSE,
    SIBLING,
    
    // Extended Family
    GRANDPARENT,
    GRANDCHILD,
    UNCLE,
    AUNT,
    NEPHEW,
    NIECE,
    COUSIN,
    
    // In-laws
    FATHER_IN_LAW,
    MOTHER_IN_LAW,
    SON_IN_LAW,
    DAUGHTER_IN_LAW,
    BROTHER_IN_LAW,
    SISTER_IN_LAW,
    
    // Indian Specific Relations
    MATERNAL_UNCLE,
    PATERNAL_UNCLE,
    MATERNAL_AUNT,
    PATERNAL_AUNT,
    MATERNAL_GRANDFATHER,
    PATERNAL_GRANDFATHER,
    MATERNAL_GRANDMOTHER,
    PATERNAL_GRANDMOTHER
}
