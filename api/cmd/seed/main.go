package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"employee-directory-api/internal/auth"
	"employee-directory-api/internal/config"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var firstNames = []string{"Ada", "Grace", "Alan", "Katherine", "Linus", "Marie", "Nikola", "Ada", "John", "Jane", "Alex", "Sam", "Taylor", "Jordan", "Morgan", "Casey", "Riley", "Avery", "Quinn", "Blake", "Drew", "Emery", "Finley", "Hayden", "Jessie", "Kai", "Logan", "Parker", "Reese", "Sage", "Rowan", "Shawn", "Toni", "Val", "Ari", "Cam", "Dana", "Emil", "Frank", "Gale", "Harley", "Ira", "Jules", "Kelly", "Leslie", "Mick", "Noel", "Pat", "Robin", "Sidney", "Tracy", "Whitney", "Aiden", "Bella", "Carlos", "Diana", "Ethan", "Fiona", "George", "Hannah", "Ivan", "Julia", "Kevin", "Laura", "Mason", "Nora", "Oscar", "Paula", "Quentin", "Rachel", "Steve", "Tina", "Uma", "Victor", "Wendy", "Xavier", "Yara", "Zane"}
var lastNames = []string{"Lovelace", "Hopper", "Turing", "Johnson", "Torvalds", "Curie", "Tesla", "Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis", "Rodriguez", "Martinez", "Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson", "Thomas", "Taylor", "Moore", "Jackson", "Martin", "Lee", "Perez", "Thompson", "White", "Harris", "Sanchez", "Clark", "Ramirez", "Lewis", "Robinson", "Walker", "Young", "Allen", "King", "Wright", "Scott", "Torres", "Nguyen", "Hill", "Flores", "Green", "Adams", "Nelson", "Baker", "Hall", "Rivera", "Campbell", "Mitchell", "Carter", "Roberts", "Gomez", "Phillips", "Evans", "Turner", "Diaz", "Parker", "Cruz", "Edwards", "Collins", "Reyes", "Stewart", "Morris", "Morales", "Murphy", "Cook", "Rogers", "Gutierrez", "Ortiz", "Morgan", "Cooper", "Peterson", "Bailey", "Reed", "Kelly", "Howard", "Ramos", "Kim", "Cox", "Ward", "Richardson", "Watson", "Brooks", "Chavez", "Wood", "James", "Bennett", "Gray", "Mendoza", "Ruiz", "Hughes", "Price", "Alvarez", "Castillo", "Sanders", "Patel", "Myers", "Long", "Ross", "Foster", "Jimenez", "Powell", "Jenkins", "Perry", "Russell", "Sullivan", "Bell", "Coleman", "Butler", "Henderson", "Barnes", "Gonzales", "Fisher", "Vasquez", "Jenkins"}
var departments = []string{"Engineering", "Data Science", "Finance", "HR", "Marketing", "Sales", "Support", "Operations", "Design", "Product", "Legal", "Research"}
var teams = []string{"Platform", "Infrastructure", "Frontend", "Backend", "Mobile", "AI", "Analytics", "Growth", "Brand"}
var locations = []string{"Chelsea", "Deer Island", "Remote", "NYC", "London", "Berlin", "Tokyo"}
var jobTitles = []string{"Engineer", "Senior Engineer", "Staff Engineer", "Principal Engineer", "Engineering Manager", "Product Manager", "Designer", "Data Scientist", "Analyst", "HR Manager", "Sales Rep"}

func main() {
	dsn := config.Load().DatabaseURL
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	// get default org
	var orgID uuid.UUID
	err = pool.QueryRow(ctx, `SELECT id FROM organizations LIMIT 1`).Scan(&orgID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("org: %s\n", orgID)

	// clear existing persons (keep org)
	mustExec(ctx, pool, `TRUNCATE persons CASCADE`)
	// re-insert org if needed (cascade deleted views? no)
	// need to re-ensure org exists
	mustExec(ctx, pool, `INSERT INTO organizations (id, name, type, country, timezone) VALUES ($1,'Default Org','company','USA','America/New_York') ON CONFLICT (id) DO NOTHING`, orgID)

	rand.Seed(time.Now().UnixNano())
	batch := 1000
	for i := 0; i < batch; i++ {
		fn := firstNames[rand.Intn(len(firstNames))]
		ln := lastNames[rand.Intn(len(lastNames))]
		dept := departments[rand.Intn(len(departments))]
		team := teams[rand.Intn(len(teams))]
		loc := locations[rand.Intn(len(locations))]
		title := jobTitles[rand.Intn(len(jobTitles))]
		personID := uuid.New()
		orgEmail := fmt.Sprintf("%s.%s%d@example.com", fn, ln, i)
		city := loc
		if loc == "Remote" {
			city = "Remote"
		}
		// person
		_, err = pool.Exec(ctx, `INSERT INTO persons (id, org_id, first_name, last_name, preferred_name, date_of_birth, gender, profile_photo_url, personal_email, org_email, phone_primary, address_line_1, city, country, is_international, source, notes, tags) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)`,
			personID, orgID, fn, ln, nil, time.Date(1980+rand.Intn(25), time.Month(1+rand.Intn(12)), 1+rand.Intn(28), 0, 0, 0, 0, time.UTC), []string{"male", "female", "non-binary"}[rand.Intn(3)], nil, fmt.Sprintf("%s.%s.personal%d@example.com", fn, ln, i), orgEmail, fmt.Sprintf("617-555-%04d", 1000+i), fmt.Sprintf("%d Main St", 100+i), city, "USA", false, []string{"referral", "job board", "direct application"}[rand.Intn(3)], nil, []string{[]string{"leadership-track", "remote", "new-hire"}[rand.Intn(3)]})
		if err != nil {
			log.Fatalf("person %d: %v", i, err)
		}

		// employment
		empID := uuid.New()
		hireDate := time.Date(2018+rand.Intn(7), time.Month(1+rand.Intn(12)), 1+rand.Intn(28), 0, 0, 0, 0, time.UTC)
		validFrom := hireDate
		status := []string{"full-time", "part-time", "contract"}[rand.Intn(3)]
		employmentType := []string{"permanent", "fixed-term"}[rand.Intn(2)]
		arrangement := []string{"on-site", "remote", "hybrid"}[rand.Intn(3)]
		_, err = pool.Exec(ctx, `INSERT INTO employment_records (id, person_id, org_id, employee_id, job_title, job_level, employment_status, employment_type, work_arrangement, department, team, office_location, desk_number, salary_amount, salary_currency, pay_frequency, hire_date, valid_from, is_current) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,true)`,
			empID, personID, orgID, fmt.Sprintf("EMP-%05d", i+1), title, []string{"L3", "L4", "Senior", "Staff"}[rand.Intn(4)], status, employmentType, arrangement, dept, team, loc, fmt.Sprintf("D-%d", rand.Intn(200)), int64(70000+rand.Intn(80000))*100, "USD", "monthly", hireDate, validFrom)
		if err != nil {
			log.Fatalf("emp %d: %v", i, err)
		}
		mustExec(ctx, pool, `INSERT INTO status_change_events (id, person_id, org_id, event_type, context, to_status, to_department, effective_date, linked_record_id, linked_record_type) VALUES ($1,$2,$3,'HIRED','employment',$4,$5,$6,$7,'employment_record')`, uuid.New(), personID, orgID, status, dept, hireDate, empID)

		// 30% get emergency contact
		if rand.Float32() < 0.3 {
			ecID := uuid.New()
			mustExec(ctx, pool, `INSERT INTO emergency_contacts (id, person_id, name, phone, email, relationship) VALUES ($1,$2,$3,$4,$5,$6)`,
				ecID, personID, fmt.Sprintf("%s %s", firstNames[rand.Intn(len(firstNames))], lastNames[rand.Intn(len(lastNames))]), fmt.Sprintf("617-555-%04d", 9000+i), fmt.Sprintf("emergency%d@example.com", i), []string{"spouse", "parent", "sibling", "friend"}[rand.Intn(4)])
		}
		if (i+1)%200 == 0 {
			fmt.Printf("seeded %d/%d\n", i+1, batch)
		}
	}
	// seed 3 known deterministic for testing
	known := []struct{ fn, ln, dept, email string }{
		{"Ada", "Lovelace", "Engineering", "ada.lovelace@example.com"},
		{"Grace", "Hopper", "Engineering", "grace.hopper@example.com"},
		{"Alan", "Turing", "Data Science", "alan.turing@example.com"},
	}
	var adminPersonID uuid.UUID
	for _, k := range known {
		pid := uuid.New()
		empID := uuid.New()
		hireDate := time.Date(2019, 4, 12, 0, 0, 0, 0, time.UTC)
		mustExec(ctx, pool, `INSERT INTO persons (id, org_id, first_name, last_name, org_email, phone_primary, city, country, tags) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`, pid, orgID, k.fn, k.ln, k.email, "617-555-0101", "Chelsea", "USA", []string{"seed"})
		mustExec(ctx, pool, `INSERT INTO employment_records (id, person_id, org_id, job_title, department, team, office_location, hire_date, valid_from, is_current) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,true)`, empID, pid, orgID, "Principal Engineer", k.dept, "Platform", "Chelsea", hireDate, hireDate)
		mustExec(ctx, pool, `INSERT INTO status_change_events (id, person_id, org_id, event_type, context, to_status, to_department, effective_date, linked_record_id, linked_record_type) VALUES ($1,$2,$3,'HIRED','employment','full-time',$4,$5,$6,'employment_record')`, uuid.New(), pid, orgID, k.dept, hireDate, empID)
		if k.fn == "Ada" {
			adminPersonID = pid
		}
	}
	// seed a default admin account bound to Ada Lovelace (username: admin / password: admin123)
	if adminPersonID != uuid.Nil {
		hash, err := auth.HashPassword("admin123")
		if err != nil {
			log.Fatal(err)
		}
		mustExec(ctx, pool, `INSERT INTO person_accounts (person_id, username, password_hash, role) VALUES ($1,'admin',$2,'admin') ON CONFLICT (username) DO UPDATE SET password_hash=EXCLUDED.password_hash, role=EXCLUDED.role, is_active=true`, adminPersonID, hash)
		fmt.Println("seeded admin account: admin / admin123")
	}
	mustExec(ctx, pool, `INSERT INTO headcount_snapshots (org_id, snapshot_date, department, location, active_count, new_entries, exits)
		SELECT org_id, CURRENT_DATE, department, office_location, COUNT(*), 0, 0
		FROM active_employees GROUP BY org_id, department, office_location
		ON CONFLICT (org_id, snapshot_date, department, location) DO UPDATE SET active_count=EXCLUDED.active_count`)
	fmt.Printf("done: seeded %d + 3 known\n", batch)
}

func mustExec(ctx context.Context, pool *pgxpool.Pool, query string, args ...any) {
	if _, err := pool.Exec(ctx, query, args...); err != nil {
		log.Fatal(err)
	}
}
