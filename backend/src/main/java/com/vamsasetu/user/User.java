package com.vamsasetu.user;

import jakarta.persistence.*;
import java.time.LocalDateTime;
import java.util.List;

@Entity
@Table(name = "users")
public class User {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    @Column(unique = true)
    private String email;
    private String name;
    private String password; // hashed
    private String roles; // admin/member/guest (CSV or join table)
    private LocalDateTime createdAt;
    private String status;
    @ElementCollection
    private List<String> familyIds; // Multiple family trees
    // getters and setters
}
