package com.vamsasetu.model;

import jakarta.persistence.*;
import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;
import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;
import org.springframework.data.jpa.domain.support.AuditingEntityListener;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.util.HashSet;
import java.util.Set;

@Entity
@Table(name = "family_members")
@EntityListeners(AuditingEntityListener.class)
public class FamilyMember {
    
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    @NotBlank
    @Size(max = 50)
    @Column(name = "first_name")
    private String firstName;
    
    @Size(max = 50)
    @Column(name = "middle_name")
    private String middleName;
    
    @NotBlank
    @Size(max = 50)
    @Column(name = "last_name")
    private String lastName;
    
    @Column(name = "date_of_birth")
    private LocalDate dateOfBirth;
    
    @Column(name = "date_of_death")
    private LocalDate dateOfDeath;
    
    @Enumerated(EnumType.STRING)
    @Column(length = 10)
    private Gender gender;
    
    @Column(name = "profile_picture")
    private String profilePicture;
    
    @Size(max = 20)
    private String phone;
    
    @Email
    @Size(max = 100)
    private String email;
    
    @Size(max = 200)
    private String occupation;
    
    @Column(columnDefinition = "TEXT")
    private String notes;
    
    @Column(name = "is_alive")
    private Boolean isAlive = true;
    
    // Address fields
    @Size(max = 200)
    private String street;
    
    @Size(max = 100)
    private String city;
    
    @Size(max = 100)
    private String state;
    
    @Size(max = 100)
    private String country;
    
    @Size(max = 20)
    @Column(name = "zip_code")
    private String zipCode;
    
    @CreatedDate
    @Column(name = "created_at", nullable = false, updatable = false)
    private LocalDateTime createdAt;
    
    @LastModifiedDate
    @Column(name = "updated_at")
    private LocalDateTime updatedAt;
    
    // Relationships
    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "user_id")
    private User user;
    
    @OneToMany(mappedBy = "fromMember", cascade = CascadeType.ALL, fetch = FetchType.LAZY)
    private Set<Relationship> outgoingRelationships = new HashSet<>();
    
    @OneToMany(mappedBy = "toMember", cascade = CascadeType.ALL, fetch = FetchType.LAZY)
    private Set<Relationship> incomingRelationships = new HashSet<>();
    
    @OneToMany(mappedBy = "member", cascade = CascadeType.ALL, fetch = FetchType.LAZY)
    private Set<EventAttendee> eventAttendances = new HashSet<>();
    
    // Constructors
    public FamilyMember() {}
    
    public FamilyMember(String firstName, String lastName, Gender gender) {
        this.firstName = firstName;
        this.lastName = lastName;
        this.gender = gender;
    }
    
    // Getters and Setters
    public Long getId() {
        return id;
    }
    
    public void setId(Long id) {
        this.id = id;
    }
    
    public String getFirstName() {
        return firstName;
    }
    
    public void setFirstName(String firstName) {
        this.firstName = firstName;
    }
    
    public String getMiddleName() {
        return middleName;
    }
    
    public void setMiddleName(String middleName) {
        this.middleName = middleName;
    }
    
    public String getLastName() {
        return lastName;
    }
    
    public void setLastName(String lastName) {
        this.lastName = lastName;
    }
    
    public LocalDate getDateOfBirth() {
        return dateOfBirth;
    }
    
    public void setDateOfBirth(LocalDate dateOfBirth) {
        this.dateOfBirth = dateOfBirth;
    }
    
    public LocalDate getDateOfDeath() {
        return dateOfDeath;
    }
    
    public void setDateOfDeath(LocalDate dateOfDeath) {
        this.dateOfDeath = dateOfDeath;
    }
    
    public Gender getGender() {
        return gender;
    }
    
    public void setGender(Gender gender) {
        this.gender = gender;
    }
    
    public String getProfilePicture() {
        return profilePicture;
    }
    
    public void setProfilePicture(String profilePicture) {
        this.profilePicture = profilePicture;
    }
    
    public String getPhone() {
        return phone;
    }
    
    public void setPhone(String phone) {
        this.phone = phone;
    }
    
    public String getEmail() {
        return email;
    }
    
    public void setEmail(String email) {
        this.email = email;
    }
    
    public String getOccupation() {
        return occupation;
    }
    
    public void setOccupation(String occupation) {
        this.occupation = occupation;
    }
    
    public String getNotes() {
        return notes;
    }
    
    public void setNotes(String notes) {
        this.notes = notes;
    }
    
    public Boolean getIsAlive() {
        return isAlive;
    }
    
    public void setIsAlive(Boolean isAlive) {
        this.isAlive = isAlive;
    }
    
    public String getStreet() {
        return street;
    }
    
    public void setStreet(String street) {
        this.street = street;
    }
    
    public String getCity() {
        return city;
    }
    
    public void setCity(String city) {
        this.city = city;
    }
    
    public String getState() {
        return state;
    }
    
    public void setState(String state) {
        this.state = state;
    }
    
    public String getCountry() {
        return country;
    }
    
    public void setCountry(String country) {
        this.country = country;
    }
    
    public String getZipCode() {
        return zipCode;
    }
    
    public void setZipCode(String zipCode) {
        this.zipCode = zipCode;
    }
    
    public LocalDateTime getCreatedAt() {
        return createdAt;
    }
    
    public void setCreatedAt(LocalDateTime createdAt) {
        this.createdAt = createdAt;
    }
    
    public LocalDateTime getUpdatedAt() {
        return updatedAt;
    }
    
    public void setUpdatedAt(LocalDateTime updatedAt) {
        this.updatedAt = updatedAt;
    }
    
    public User getUser() {
        return user;
    }
    
    public void setUser(User user) {
        this.user = user;
    }
    
    public Set<Relationship> getOutgoingRelationships() {
        return outgoingRelationships;
    }
    
    public void setOutgoingRelationships(Set<Relationship> outgoingRelationships) {
        this.outgoingRelationships = outgoingRelationships;
    }
    
    public Set<Relationship> getIncomingRelationships() {
        return incomingRelationships;
    }
    
    public void setIncomingRelationships(Set<Relationship> incomingRelationships) {
        this.incomingRelationships = incomingRelationships;
    }
    
    public Set<EventAttendee> getEventAttendances() {
        return eventAttendances;
    }
    
    public void setEventAttendances(Set<EventAttendee> eventAttendances) {
        this.eventAttendances = eventAttendances;
    }
}
