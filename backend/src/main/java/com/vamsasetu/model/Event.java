package com.vamsasetu.model;

import jakarta.persistence.*;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;
import org.springframework.data.jpa.domain.support.AuditingEntityListener;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.HashSet;
import java.util.Set;

@Entity
@Table(name = "events")
@EntityListeners(AuditingEntityListener.class)
public class Event {
    
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    @NotBlank
    @Size(max = 200)
    private String title;
    
    @Column(columnDefinition = "TEXT")
    private String description;
    
    @Enumerated(EnumType.STRING)
    @Column(name = "event_type", nullable = false)
    private EventType eventType;
    
    @NotNull
    @Column(name = "event_date", nullable = false)
    private LocalDate eventDate;
    
    @Column(name = "event_time")
    private LocalTime eventTime;
    
    @Size(max = 200)
    private String location;
    
    @Column(name = "is_recurring")
    private Boolean isRecurring = false;
    
    @Column(name = "recurrence_frequency")
    @Enumerated(EnumType.STRING)
    private RecurrenceFrequency recurrenceFrequency;
    
    @Column(name = "recurrence_interval")
    private Integer recurrenceInterval = 1;
    
    @Column(name = "recurrence_end_date")
    private LocalDate recurrenceEndDate;
    
    @Column(name = "recurrence_occurrences")
    private Integer recurrenceOccurrences;
    
    // Reminder Settings
    @Column(name = "reminder_enabled")
    private Boolean reminderEnabled = true;
    
    @Column(name = "reminder_advance_time")
    private Integer reminderAdvanceTime = 24; // in hours
    
    @Column(name = "reminder_email")
    private Boolean reminderEmail = true;
    
    @Column(name = "reminder_sms")
    private Boolean reminderSms = false;
    
    @Column(name = "reminder_whatsapp")
    private Boolean reminderWhatsapp = false;
    
    @Column(name = "reminder_push")
    private Boolean reminderPush = true;
    
    @Column(name = "custom_email_message", columnDefinition = "TEXT")
    private String customEmailMessage;
    
    @Column(name = "custom_sms_message", columnDefinition = "TEXT")
    private String customSmsMessage;
    
    @Column(name = "custom_whatsapp_message", columnDefinition = "TEXT")
    private String customWhatsappMessage;
    
    @CreatedDate
    @Column(name = "created_at", nullable = false, updatable = false)
    private LocalDateTime createdAt;
    
    @LastModifiedDate
    @Column(name = "updated_at")
    private LocalDateTime updatedAt;
    
    // Relationships
    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "created_by", nullable = false)
    @NotNull
    private User createdBy;
    
    @OneToMany(mappedBy = "event", cascade = CascadeType.ALL, fetch = FetchType.LAZY)
    private Set<EventAttendee> attendees = new HashSet<>();
    
    // Constructors
    public Event() {}
    
    public Event(String title, EventType eventType, LocalDate eventDate, User createdBy) {
        this.title = title;
        this.eventType = eventType;
        this.eventDate = eventDate;
        this.createdBy = createdBy;
    }
    
    // Getters and Setters
    public Long getId() {
        return id;
    }
    
    public void setId(Long id) {
        this.id = id;
    }
    
    public String getTitle() {
        return title;
    }
    
    public void setTitle(String title) {
        this.title = title;
    }
    
    public String getDescription() {
        return description;
    }
    
    public void setDescription(String description) {
        this.description = description;
    }
    
    public EventType getEventType() {
        return eventType;
    }
    
    public void setEventType(EventType eventType) {
        this.eventType = eventType;
    }
    
    public LocalDate getEventDate() {
        return eventDate;
    }
    
    public void setEventDate(LocalDate eventDate) {
        this.eventDate = eventDate;
    }
    
    public LocalTime getEventTime() {
        return eventTime;
    }
    
    public void setEventTime(LocalTime eventTime) {
        this.eventTime = eventTime;
    }
    
    public String getLocation() {
        return location;
    }
    
    public void setLocation(String location) {
        this.location = location;
    }
    
    public Boolean getIsRecurring() {
        return isRecurring;
    }
    
    public void setIsRecurring(Boolean isRecurring) {
        this.isRecurring = isRecurring;
    }
    
    public RecurrenceFrequency getRecurrenceFrequency() {
        return recurrenceFrequency;
    }
    
    public void setRecurrenceFrequency(RecurrenceFrequency recurrenceFrequency) {
        this.recurrenceFrequency = recurrenceFrequency;
    }
    
    public Integer getRecurrenceInterval() {
        return recurrenceInterval;
    }
    
    public void setRecurrenceInterval(Integer recurrenceInterval) {
        this.recurrenceInterval = recurrenceInterval;
    }
    
    public LocalDate getRecurrenceEndDate() {
        return recurrenceEndDate;
    }
    
    public void setRecurrenceEndDate(LocalDate recurrenceEndDate) {
        this.recurrenceEndDate = recurrenceEndDate;
    }
    
    public Integer getRecurrenceOccurrences() {
        return recurrenceOccurrences;
    }
    
    public void setRecurrenceOccurrences(Integer recurrenceOccurrences) {
        this.recurrenceOccurrences = recurrenceOccurrences;
    }
    
    public Boolean getReminderEnabled() {
        return reminderEnabled;
    }
    
    public void setReminderEnabled(Boolean reminderEnabled) {
        this.reminderEnabled = reminderEnabled;
    }
    
    public Integer getReminderAdvanceTime() {
        return reminderAdvanceTime;
    }
    
    public void setReminderAdvanceTime(Integer reminderAdvanceTime) {
        this.reminderAdvanceTime = reminderAdvanceTime;
    }
    
    public Boolean getReminderEmail() {
        return reminderEmail;
    }
    
    public void setReminderEmail(Boolean reminderEmail) {
        this.reminderEmail = reminderEmail;
    }
    
    public Boolean getReminderSms() {
        return reminderSms;
    }
    
    public void setReminderSms(Boolean reminderSms) {
        this.reminderSms = reminderSms;
    }
    
    public Boolean getReminderWhatsapp() {
        return reminderWhatsapp;
    }
    
    public void setReminderWhatsapp(Boolean reminderWhatsapp) {
        this.reminderWhatsapp = reminderWhatsapp;
    }
    
    public Boolean getReminderPush() {
        return reminderPush;
    }
    
    public void setReminderPush(Boolean reminderPush) {
        this.reminderPush = reminderPush;
    }
    
    public String getCustomEmailMessage() {
        return customEmailMessage;
    }
    
    public void setCustomEmailMessage(String customEmailMessage) {
        this.customEmailMessage = customEmailMessage;
    }
    
    public String getCustomSmsMessage() {
        return customSmsMessage;
    }
    
    public void setCustomSmsMessage(String customSmsMessage) {
        this.customSmsMessage = customSmsMessage;
    }
    
    public String getCustomWhatsappMessage() {
        return customWhatsappMessage;
    }
    
    public void setCustomWhatsappMessage(String customWhatsappMessage) {
        this.customWhatsappMessage = customWhatsappMessage;
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
    
    public User getCreatedBy() {
        return createdBy;
    }
    
    public void setCreatedBy(User createdBy) {
        this.createdBy = createdBy;
    }
    
    public Set<EventAttendee> getAttendees() {
        return attendees;
    }
    
    public void setAttendees(Set<EventAttendee> attendees) {
        this.attendees = attendees;
    }
}
