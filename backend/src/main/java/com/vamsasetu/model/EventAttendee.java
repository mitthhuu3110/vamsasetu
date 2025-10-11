package com.vamsasetu.model;

import jakarta.persistence.*;
import jakarta.validation.constraints.NotNull;
import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;
import org.springframework.data.jpa.domain.support.AuditingEntityListener;

import java.time.LocalDateTime;

@Entity
@Table(name = "event_attendees")
@EntityListeners(AuditingEntityListener.class)
public class EventAttendee {
    
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "event_id", nullable = false)
    @NotNull
    private Event event;
    
    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "member_id", nullable = false)
    @NotNull
    private FamilyMember member;
    
    @Enumerated(EnumType.STRING)
    @Column(name = "attendee_role", nullable = false)
    private AttendeeRole attendeeRole = AttendeeRole.GUEST;
    
    @Enumerated(EnumType.STRING)
    @Column(name = "rsvp_status", nullable = false)
    private RsvpStatus rsvpStatus = RsvpStatus.PENDING;
    
    @CreatedDate
    @Column(name = "created_at", nullable = false, updatable = false)
    private LocalDateTime createdAt;
    
    @LastModifiedDate
    @Column(name = "updated_at")
    private LocalDateTime updatedAt;
    
    // Constructors
    public EventAttendee() {}
    
    public EventAttendee(Event event, FamilyMember member, AttendeeRole attendeeRole) {
        this.event = event;
        this.member = member;
        this.attendeeRole = attendeeRole;
    }
    
    // Getters and Setters
    public Long getId() {
        return id;
    }
    
    public void setId(Long id) {
        this.id = id;
    }
    
    public Event getEvent() {
        return event;
    }
    
    public void setEvent(Event event) {
        this.event = event;
    }
    
    public FamilyMember getMember() {
        return member;
    }
    
    public void setMember(FamilyMember member) {
        this.member = member;
    }
    
    public AttendeeRole getAttendeeRole() {
        return attendeeRole;
    }
    
    public void setAttendeeRole(AttendeeRole attendeeRole) {
        this.attendeeRole = attendeeRole;
    }
    
    public RsvpStatus getRsvpStatus() {
        return rsvpStatus;
    }
    
    public void setRsvpStatus(RsvpStatus rsvpStatus) {
        this.rsvpStatus = rsvpStatus;
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
}
