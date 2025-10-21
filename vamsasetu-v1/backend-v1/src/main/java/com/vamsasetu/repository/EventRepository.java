package com.vamsasetu.repository;

import com.vamsasetu.model.Event;
import com.vamsasetu.model.EventType;
import com.vamsasetu.model.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.time.LocalDate;
import java.util.List;

@Repository
public interface EventRepository extends JpaRepository<Event, Long> {
    
    List<Event> findByCreatedBy(User createdBy);
    
    List<Event> findByCreatedByOrderByEventDateAsc(User createdBy);
    
    List<Event> findByEventType(EventType eventType);
    
    List<Event> findByEventDate(LocalDate eventDate);
    
    List<Event> findByEventDateBetween(LocalDate startDate, LocalDate endDate);
    
    @Query("SELECT e FROM Event e WHERE e.createdBy = :user AND e.eventDate >= :date ORDER BY e.eventDate ASC")
    List<Event> findUpcomingEventsByUser(@Param("user") User user, @Param("date") LocalDate date);
    
    @Query("SELECT e FROM Event e WHERE e.createdBy = :user AND e.eventDate < :date ORDER BY e.eventDate DESC")
    List<Event> findPastEventsByUser(@Param("user") User user, @Param("date") LocalDate date);
    
    @Query("SELECT e FROM Event e WHERE e.createdBy = :user AND e.eventType = :eventType ORDER BY e.eventDate ASC")
    List<Event> findByUserAndEventType(@Param("user") User user, @Param("eventType") EventType eventType);
    
    @Query("SELECT e FROM Event e WHERE e.createdBy = :user AND " +
           "(LOWER(e.title) LIKE LOWER(CONCAT('%', :query, '%')) OR " +
           "LOWER(e.description) LIKE LOWER(CONCAT('%', :query, '%')))")
    List<Event> findByUserAndTitleOrDescriptionContaining(@Param("user") User user, @Param("query") String query);
    
    @Query("SELECT e FROM Event e WHERE e.createdBy = :user AND e.isRecurring = true")
    List<Event> findRecurringEventsByUser(@Param("user") User user);
    
    @Query("SELECT e FROM Event e WHERE e.createdBy = :user AND e.reminderEnabled = true AND " +
           "e.eventDate BETWEEN :startDate AND :endDate")
    List<Event> findEventsWithRemindersInDateRange(@Param("user") User user, 
                                                   @Param("startDate") LocalDate startDate, 
                                                   @Param("endDate") LocalDate endDate);
}
