package main

import (
	"context"
	"log"
	"time"

	"vamsasetu/backend/internal/config"
	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/internal/repository"
	"vamsasetu/backend/pkg/neo4j"
	"vamsasetu/backend/pkg/postgres"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connections
	pgClient, err := postgres.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pgClient.Close()

	neo4jClient, err := neo4j.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err)
	}
	defer neo4jClient.Close(context.Background())

	ctx := context.Background()

	// Initialize repositories
	userRepo := repository.NewUserRepository(pgClient.DB)
	memberRepo := repository.NewMemberRepository(neo4jClient)
	relationshipRepo := repository.NewRelationshipRepository(neo4jClient)
	eventRepo := repository.NewEventRepository(pgClient.DB)

	log.Println("Starting seed data creation...")

	// 1. Create demo user
	log.Println("Creating demo user...")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Demo@1234"), bcrypt.DefaultCost)
	demoUser := &models.User{
		Email:        "demo@vamsasetu.com",
		PasswordHash: string(hashedPassword),
		Name:         "Demo User",
		Role:         "owner",
	}
	if err := userRepo.Create(ctx, demoUser); err != nil {
		log.Printf("Warning: Demo user might already exist: %v", err)
	} else {
		log.Printf("Created demo user: %s", demoUser.Email)
	}

	// 2. Create 12 members across 3 generations
	log.Println("Creating family members...")
	
	// Generation 1 (Grandparents)
	grandfather := createMember(ctx, memberRepo, "Ramesh Kumar", "1950-03-15", "male", "ramesh@example.com", "+91-9876543210")
	grandmother := createMember(ctx, memberRepo, "Lakshmi Devi", "1952-07-20", "female", "lakshmi@example.com", "+91-9876543211")
	
	// Generation 2 (Parents and Uncles/Aunts)
	father := createMember(ctx, memberRepo, "Suresh Kumar", "1975-05-10", "male", "suresh@example.com", "+91-9876543212")
	mother := createMember(ctx, memberRepo, "Priya Kumari", "1977-08-25", "female", "priya@example.com", "+91-9876543213")
	uncle := createMember(ctx, memberRepo, "Mahesh Kumar", "1978-11-30", "male", "mahesh@example.com", "+91-9876543214")
	aunt := createMember(ctx, memberRepo, "Radha Devi", "1980-02-14", "female", "radha@example.com", "+91-9876543215")
	
	// Generation 3 (Children and Cousins)
	son1 := createMember(ctx, memberRepo, "Arun Kumar", "2000-01-15", "male", "arun@example.com", "+91-9876543216")
	daughter1 := createMember(ctx, memberRepo, "Anjali Kumari", "2002-06-20", "female", "anjali@example.com", "+91-9876543217")
	son2 := createMember(ctx, memberRepo, "Karthik Kumar", "2003-09-10", "male", "karthik@example.com", "+91-9876543218")
	daughter2 := createMember(ctx, memberRepo, "Divya Kumari", "2005-12-05", "female", "divya@example.com", "+91-9876543219")
	cousin1 := createMember(ctx, memberRepo, "Ravi Kumar", "2001-03-22", "male", "ravi@example.com", "+91-9876543220")
	cousin2 := createMember(ctx, memberRepo, "Meera Kumari", "2004-07-18", "female", "meera@example.com", "+91-9876543221")

	// 3. Create relationships
	log.Println("Creating relationships...")
	
	// Generation 1 relationships
	createRelationship(ctx, relationshipRepo, grandfather.ID, grandmother.ID, "SPOUSE_OF")
	
	// Generation 2 relationships
	createRelationship(ctx, relationshipRepo, grandfather.ID, father.ID, "PARENT_OF")
	createRelationship(ctx, relationshipRepo, grandmother.ID, father.ID, "PARENT_OF")
	createRelationship(ctx, relationshipRepo, grandfather.ID, uncle.ID, "PARENT_OF")
	createRelationship(ctx, relationshipRepo, grandmother.ID, uncle.ID, "PARENT_OF")
	createRelationship(ctx, relationshipRepo, father.ID, uncle.ID, "SIBLING_OF")
	createRelationship(ctx, relationshipRepo, father.ID, mother.ID, "SPOUSE_OF")
	createRelationship(ctx, relationshipRepo, uncle.ID, aunt.ID, "SPOUSE_OF")
	
	// Generation 3 relationships
	createRelationship(ctx, relationshipRepo, father.ID, son1.ID, "PARENT_OF")
	createRelationship(ctx, relationshipRepo, mother.ID, son1.ID, "PARENT_OF")
	createRelationship(ctx, relationshipRepo, father.ID, daughter1.ID, "PARENT_OF")
	createRelationship(ctx, relationshipRepo, mother.ID, daughter1.ID, "PARENT_OF")
	createRelationship(ctx, relationshipRepo, father.ID, son2.ID, "PARENT_OF")
	createRelationship(ctx, relationshipRepo, mother.ID, son2.ID, "PARENT_OF")
	createRelationship(ctx, relationshipRepo, father.ID, daughter2.ID, "PARENT_OF")
	createRelationship(ctx, relationshipRepo, mother.ID, daughter2.ID, "PARENT_OF")
	createRelationship(ctx, relationshipRepo, uncle.ID, cousin1.ID, "PARENT_OF")
	createRelationship(ctx, relationshipRepo, aunt.ID, cousin1.ID, "PARENT_OF")
	createRelationship(ctx, relationshipRepo, uncle.ID, cousin2.ID, "PARENT_OF")
	createRelationship(ctx, relationshipRepo, aunt.ID, cousin2.ID, "PARENT_OF")
	
	// Sibling relationships
	createRelationship(ctx, relationshipRepo, son1.ID, daughter1.ID, "SIBLING_OF")
	createRelationship(ctx, relationshipRepo, son1.ID, son2.ID, "SIBLING_OF")
	createRelationship(ctx, relationshipRepo, son1.ID, daughter2.ID, "SIBLING_OF")
	createRelationship(ctx, relationshipRepo, daughter1.ID, son2.ID, "SIBLING_OF")
	createRelationship(ctx, relationshipRepo, daughter1.ID, daughter2.ID, "SIBLING_OF")
	createRelationship(ctx, relationshipRepo, son2.ID, daughter2.ID, "SIBLING_OF")
	createRelationship(ctx, relationshipRepo, cousin1.ID, cousin2.ID, "SIBLING_OF")

	// 4. Create 5 events
	log.Println("Creating events...")
	
	// Birthday events
	createEvent(ctx, eventRepo, "Arun's Birthday", "Celebrating Arun's birthday", time.Now().AddDate(0, 0, 5), "birthday", []string{son1.ID})
	createEvent(ctx, eventRepo, "Anjali's Birthday", "Celebrating Anjali's birthday", time.Now().AddDate(0, 0, 15), "birthday", []string{daughter1.ID})
	
	// Anniversary event
	createEvent(ctx, eventRepo, "Suresh & Priya Anniversary", "25th wedding anniversary celebration", time.Now().AddDate(0, 0, 30), "anniversary", []string{father.ID, mother.ID})
	
	// Ceremony event
	createEvent(ctx, eventRepo, "Family Puja", "Annual family puja ceremony", time.Now().AddDate(0, 1, 0), "ceremony", []string{grandfather.ID, grandmother.ID, father.ID, mother.ID, uncle.ID, aunt.ID})
	
	// Custom event
	createEvent(ctx, eventRepo, "Family Reunion", "Annual family get-together", time.Now().AddDate(0, 2, 0), "custom", []string{grandfather.ID, grandmother.ID, father.ID, mother.ID, uncle.ID, aunt.ID, son1.ID, daughter1.ID, son2.ID, daughter2.ID, cousin1.ID, cousin2.ID})

	log.Println("Seed data creation completed successfully!")
	log.Println("Demo user credentials:")
	log.Println("  Email: demo@vamsasetu.com")
	log.Println("  Password: Demo@1234")
}

func createMember(ctx context.Context, repo *repository.MemberRepository, name, dob, gender, email, phone string) *models.Member {
	dobTime, _ := time.Parse("2006-01-02", dob)
	member := &models.Member{
		Name:        name,
		DateOfBirth: dobTime,
		Gender:      gender,
		Email:       email,
		Phone:       phone,
	}
	if err := repo.Create(ctx, member); err != nil {
		log.Printf("Warning: Failed to create member %s: %v", name, err)
	} else {
		log.Printf("Created member: %s (ID: %s)", name, member.ID)
	}
	return member
}

func createRelationship(ctx context.Context, repo *repository.RelationshipRepository, fromID, toID, relType string) {
	rel := &models.Relationship{
		FromID: fromID,
		ToID:   toID,
		Type:   relType,
	}
	if err := repo.Create(ctx, rel); err != nil {
		log.Printf("Warning: Failed to create relationship %s -> %s (%s): %v", fromID, toID, relType, err)
	} else {
		log.Printf("Created relationship: %s -> %s (%s)", fromID, toID, relType)
	}
}

func createEvent(ctx context.Context, repo *repository.EventRepository, title, description string, eventDate time.Time, eventType string, memberIDs []string) {
	event := &models.Event{
		Title:       title,
		Description: description,
		EventDate:   eventDate,
		EventType:   eventType,
		CreatedBy:   1, // Demo user ID
	}
	event.SetMemberIDs(memberIDs)
	if err := repo.Create(ctx, event); err != nil {
		log.Printf("Warning: Failed to create event %s: %v", title, err)
	} else {
		log.Printf("Created event: %s (Date: %s)", title, eventDate.Format("2006-01-02"))
	}
}
