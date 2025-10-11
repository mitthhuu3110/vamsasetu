package com.vamsasetu.repository;

import com.vamsasetu.model.FamilyMember;
import com.vamsasetu.model.Gender;
import com.vamsasetu.model.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;

@Repository
public interface FamilyMemberRepository extends JpaRepository<FamilyMember, Long> {
    
    List<FamilyMember> findByUser(User user);
    
    @Query("SELECT fm FROM FamilyMember fm WHERE fm.user = :user ORDER BY fm.firstName, fm.lastName")
    List<FamilyMember> findByUserOrderByName(@Param("user") User user);
    
    @Query("SELECT fm FROM FamilyMember fm WHERE fm.user = :user AND fm.isAlive = :isAlive")
    List<FamilyMember> findByUserAndIsAlive(@Param("user") User user, @Param("isAlive") Boolean isAlive);
    
    @Query("SELECT fm FROM FamilyMember fm WHERE fm.user = :user AND fm.gender = :gender")
    List<FamilyMember> findByUserAndGender(@Param("user") User user, @Param("gender") Gender gender);
    
    @Query("SELECT fm FROM FamilyMember fm WHERE fm.user = :user AND " +
           "(LOWER(fm.firstName) LIKE LOWER(CONCAT('%', :query, '%')) OR " +
           "LOWER(fm.lastName) LIKE LOWER(CONCAT('%', :query, '%')) OR " +
           "LOWER(fm.middleName) LIKE LOWER(CONCAT('%', :query, '%')))")
    List<FamilyMember> findByUserAndNameContaining(@Param("user") User user, @Param("query") String query);
    
    @Query("SELECT fm FROM FamilyMember fm WHERE fm.user = :user AND fm.id = :id")
    Optional<FamilyMember> findByUserAndId(@Param("user") User user, @Param("id") Long id);
    
    @Query("SELECT COUNT(fm) FROM FamilyMember fm WHERE fm.user = :user")
    Long countByUser(@Param("user") User user);
    
    @Query("SELECT COUNT(fm) FROM FamilyMember fm WHERE fm.user = :user AND fm.isAlive = true")
    Long countAliveByUser(@Param("user") User user);
}
