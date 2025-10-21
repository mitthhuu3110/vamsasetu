package com.vamsasetu.event;

import jakarta.persistence.*;
import java.time.LocalDate;
import java.util.List;

@Entity
@Table(name = "events")
public class Event {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    private String name;
    private String type; // birthday, anniversary, custom
    private LocalDate eventDate;
    private String description;
    private String createdBy;
    @ElementCollection
    private List<String> memberIds; // Family members associated
    // getters and setters
}
