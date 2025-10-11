package com.vamsasetu.repository;

import com.vamsasetu.model.FamilyMember;
import com.vamsasetu.model.Relationship;
import com.vamsasetu.model.RelationshipType;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;

@Repository
public interface RelationshipRepository extends JpaRepository<Relationship, Long> {
    
    List<Relationship> findByFromMember(FamilyMember fromMember);
    
    List<Relationship> findByToMember(FamilyMember toMember);
    
    List<Relationship> findByFromMemberAndIsActive(FamilyMember fromMember, Boolean isActive);
    
    List<Relationship> findByToMemberAndIsActive(FamilyMember toMember, Boolean isActive);
    
    @Query("SELECT r FROM Relationship r WHERE (r.fromMember = :member OR r.toMember = :member) AND r.isActive = true")
    List<Relationship> findByMemberAndIsActive(@Param("member") FamilyMember member);
    
    @Query("SELECT r FROM Relationship r WHERE r.fromMember = :fromMember AND r.toMember = :toMember")
    Optional<Relationship> findByFromMemberAndToMember(@Param("fromMember") FamilyMember fromMember, 
                                                      @Param("toMember") FamilyMember toMember);
    
    @Query("SELECT r FROM Relationship r WHERE (r.fromMember = :fromMember AND r.toMember = :toMember) OR " +
           "(r.fromMember = :toMember AND r.toMember = :fromMember)")
    List<Relationship> findBidirectionalRelationship(@Param("fromMember") FamilyMember fromMember, 
                                                    @Param("toMember") FamilyMember toMember);
    
    List<Relationship> findByRelationshipType(RelationshipType relationshipType);
    
    @Query("SELECT r FROM Relationship r WHERE r.fromMember.user = :user OR r.toMember.user = :user")
    List<Relationship> findByUser(@Param("user") com.vamsasetu.model.User user);
    
    @Query("SELECT r FROM Relationship r WHERE (r.fromMember.user = :user OR r.toMember.user = :user) AND r.isActive = true")
    List<Relationship> findByUserAndIsActive(@Param("user") com.vamsasetu.model.User user);
}
