package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"thesis-management-app/handlers"
	"thesis-management-app/pkgs/ldap"
	"thesis-management-app/pkgs/logging"
	"thesis-management-app/pkgs/server"
	"thesis-management-app/types/sqlite"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func main() {

	logging.SetupLogger()

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	server.MyS.Router = chi.NewRouter()
	ldap.SetupLDAP()

	db_path := os.Getenv("DB_PATH")
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//Uncoment this to populate database with example data
	// populateDatabase(db)
	server.MyS.DB = &sqlite.Model{DB: db}

	server.MyS.DB.DB.SetConnMaxLifetime(time.Minute)
	server.MyS.DB.DB.SetMaxOpenConns(1)
	_, _ = server.MyS.DB.DB.Exec("PRAGMA busy_timeout = 5000") // Log retries
	_, err = server.MyS.DB.DB.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		log.Fatal("Failed to enable WAL mode:", err)
	}

	server.MyS.Router.Use(middleware.Logger)
	server.MyS.Router.Use(handlers.WithUser)
	server.MyS.Router.Use(handlers.RefreshSession)

	server.MyS.Router.Group(func(r chi.Router) {
		r.Handle("/*", public())
		r.Get("/", handlers.Make(handlers.HandleHome))
		r.Get("/login", handlers.Make(handlers.HandleLogin))
		r.Post("/login", handlers.Make(handlers.HandleLoginPost))
		r.Post("/logout", handlers.Make(handlers.HandleLogoutPost))
	})

	server.MyS.Router.Group(func(r chi.Router) {
		r.Use(handlers.WithAuth)
		r.Get("/ongoing/clear-new", handlers.Make(handlers.HandleOngoingClearNew))
		r.Get("/ongoing/generate_excel", handlers.Make(handlers.HandleRealizedGenerateExcel))
		r.Get("/ongoing/filter", handlers.Make(handlers.HandleOngoingFiltered))
		r.Get("/ongoing/details/{id}", handlers.Make(handlers.HandleOngoingDetails))
		r.Get("/ongoing/next_page", handlers.Make(handlers.HandleOngoingNext))
		r.Get("/ongoing/previous_page", handlers.Make(handlers.HandleOngoingPrev))
		r.Get("/ongoing/{id}", handlers.Make(handlers.HandleOngoingEntry))
		r.Get("/ongoing", handlers.Make(handlers.HandleOngoing))
		r.Get("/realized/generate_excel", handlers.Make(handlers.HandleRealizedGenerateExcel))
		r.Get("/realized/filter", handlers.Make(handlers.HandleRealizedFiltered))
		r.Get("/realized/clear-new", handlers.Make(handlers.HandleRealizedClearNew))
		r.Get("/realized/details/{id}", handlers.Make(handlers.HandleRealizedDetails))
		r.Get("/realized/next_page", handlers.Make(handlers.HandleRealizedNext))
		r.Get("/realized/previous_page", handlers.Make(handlers.HandleRealizedPrev))
		r.Get("/realized/{id}", handlers.Make(handlers.HandleRealizedEntry))
		r.Get("/realized/autocompleteThesisTitlePolish", handlers.Make(handlers.HandleAutocompleteThesisTitlePolish))
		r.Get("/realized/autocompleteStudentSurname", handlers.Make(handlers.HandleAutocompleteStudentSurname))
		r.Get("/realized/autocompleteStudentNumber", handlers.Make(handlers.HandleAutocompleteStudentNumber))
		r.Get("/realized/autocompleteStudentNameAndSurname", handlers.Make(handlers.HandleAutocompleteStudentNameAndSurname))
		r.Get("/realized/autocompleteSupervisorName", handlers.Make(handlers.HandleAutocompleteSupervisorName))
		r.Get("/realized/autocompleteSupervisorSurname", handlers.Make(handlers.HandleAutocompleteSupervisorSurname))
		r.Get("/realized/autocompleteSupervisorTitle", handlers.Make(handlers.HandleAutocompleteSupervisorTitle))
		r.Get("/realized/autocompleteSupervisorNameAndSurname", handlers.Make(handlers.HandleAutocompleteSupervisorNameAndSurname))
		r.Get("/realized/autocompleteAssistantSupervisorName", handlers.Make(handlers.HandleAutocompleteAssistantSupervisorName))
		r.Get("/realized/autocompleteAssistantSupervisorSurname", handlers.Make(handlers.HandleAutocompleteAssistantSupervisorSurname))
		r.Get("/realized/autocompleteAssistantSupervisorTitle", handlers.Make(handlers.HandleAutocompleteAssistantSupervisorTitle))
		r.Get("/realized/autocompleteAssistantSupervisorNameAndSurname", handlers.Make(handlers.HandleAutocompleteAssistantSupervisorNameAndSurname))
		r.Get("/realized/autocompleteReviewerName", handlers.Make(handlers.HandleAutocompleteReviewerName))
		r.Get("/realized/autocompleteReviewerSurname", handlers.Make(handlers.HandleAutocompleteReviewerSurname))
		r.Get("/realized/autocompleteReviewerTitle", handlers.Make(handlers.HandleAutocompleteReviewerTitle))
		r.Get("/realized/autocompleteReviewerNameAndSurname", handlers.Make(handlers.HandleAutocompleteReviewerNameAndSurname))
		r.Get("/realized/autocompleteChairName", handlers.Make(handlers.HandleAutocompleteChairName))
		r.Get("/realized/autocompleteChairSurname", handlers.Make(handlers.HandleAutocompleteChairSurname))
		r.Get("/realized/autocompleteChairTitle", handlers.Make(handlers.HandleAutocompleteChairTitle))
		r.Get("/realized/autocompleteCourse", handlers.Make(handlers.HandleAutocompleteCourse))
		r.Get("/realized", handlers.Make(handlers.HandleRealized))
		r.Get("/note/{realized_id}&{ongoing_id}&{user_id}", handlers.Make(handlers.HandleNote))
		r.Put("/realized/{id}", handlers.Make(handlers.HandleRealizedUpdate))
		r.Put("/ongoing/{id}", handlers.Make(handlers.HandleRealizedUpdate))
	})

	server.MyS.Router.Group(func(r chi.Router) {
		r.Use(handlers.WithAdminRights)
		r.Get("/settings", handlers.Make(handlers.HandleSettingsIndex))

		r.Get("/settings/employees", handlers.Make(handlers.HandleEmployees))
		r.Get("/settings/employees/new", handlers.Make(handlers.HandleEmployeesGetNew))
		r.Post("/settings/employees", handlers.Make(handlers.HandleEmployeesNew))
		r.Get("/settings/employees/{id}", handlers.Make(handlers.HandleEmployeesEntry))
		r.Get("/settings/employees/details/{id}", handlers.Make(handlers.HandleEmployeesDetails))
		r.Put("/settings/employees/{id}", handlers.Make(handlers.HandleEmployeesUpdate))
		r.Get("/settings/employees/clear-new", handlers.Make(handlers.HandleEmployeesClearNew))
		r.Get("/settings/employees/sorted", handlers.Make(handlers.HandleSortingEmps))

		r.Get("/settings/courses", handlers.Make(handlers.HandleCourses))
		r.Get("/settings/courses/new", handlers.Make(handlers.HandleCoursesGetNew))
		r.Post("/settings/courses", handlers.Make(handlers.HandleCoursesNew))
		r.Get("/settings/courses/{id}", handlers.Make(handlers.HandleCoursesEntry))
		r.Get("/settings/courses/details/{id}", handlers.Make(handlers.HandleCoursesDetails))
		r.Put("/settings/courses/{id}", handlers.Make(handlers.HandleCoursesUpdate))
		r.Get("/settings/courses/clear-new", handlers.Make(handlers.HandleCoursesClearNew))
		r.Get("/settings/courses/sorted", handlers.Make(handlers.HandleSortingCourses))

		// Routes for specializations
		r.Get("/settings/specs", handlers.Make(handlers.HandleSpecializations))
		r.Get("/settings/specs/new", handlers.Make(handlers.HandleSpecializationsGetNew))
		r.Post("/settings/specs", handlers.Make(handlers.HandleSpecializationsNew))
		r.Get("/settings/specs/{id}", handlers.Make(handlers.HandleSpecializationsEntry))
		r.Get("/settings/specs/details/{id}", handlers.Make(handlers.HandleSpecializationsDetails))
		r.Put("/settings/specs/{id}", handlers.Make(handlers.HandleSpecializationsUpdate))
		r.Get("/settings/specs/clear-new", handlers.Make(handlers.HandleSpecializationsClearNew))
		r.Get("/settings/specs/sorted", handlers.Make(handlers.HandleSortingSpecs))

		r.Get("/realized/new", handlers.Make(handlers.HandleRealizedGetNew))
		r.Post("/realized", handlers.Make(handlers.HandleRealizedNew))

		r.Get("/ongoing/new", handlers.Make(handlers.HandleOngoingGetNew))
		r.Post("/ongoing", handlers.Make(handlers.HandleOngoingNew))
	})

	listenAddr := os.Getenv("LISTEN_ADDR")

	slog.Info("HTTP server", "ListenAddr", listenAddr)
	http.ListenAndServe(listenAddr, server.MyS.Router)
}

func populateDatabase(db *sql.DB) error {
	// Queries to populate the database
	queries := []string{
		// Insert University Employees
		`INSERT INTO University_Employee (first_name, last_name, current_academic_title, department_unit) VALUES
		('Anna', 'Kowalska', 'Dr.', 'Department of Mathematics'),
		('Piotr', 'Nowak', 'Mgr. inż.', 'Department of Computer Science'),
		('Olgierd', 'Kwiatkowski', 'Mgr. inż.', 'Department of Computer Science'),
		('Jan', 'Paweł II', 'Dr. hab.', 'Department of Computer Science'),
		('Antoni', 'Ramos', 'Mgr. inż.', 'Department of Computer Science and Microelectronics'),
		('Witold', 'Chmielewski', 'Mgr. inż.', 'Department of Computer Science'),
		('Kamil', 'Olech', 'Dr. hab.', 'DMCS'),
		('Janusz', 'Wiśniewski', 'Prof. dr hab.', 'Department of Physics'),
		('Katarzyna', 'Zielińska', 'Inż.', 'Department of Chemistry'),
		('Tomasz', 'Wójcik', 'Dr. hab.', 'Department of Mechanical Engineering'),
		('Aleksandra', 'Lewandowska', 'Mgr.', 'Department of Biology'),
		('Małgorzata', 'Kaczmarek', 'Prof. dr hab.', 'Department of Environmental Science'),
		('Michał', 'Szymański', 'Dr.', 'Department of Economics'),
		('Paweł', 'Jabłoński', 'Mgr. inż.', 'Department of Electrical Engineering'),
		('Dorota', 'Kamińska', 'Inż.', 'Department of Civil Engineering'),
		('Agnieszka', 'Zawadzka', 'Dr. hab.', 'Department of Management'),
		('Krzysztof', 'Majewski', 'Prof. dr hab.', 'Department of Humanities');`,

		// Insert Students
		`INSERT INTO Student (student_number, first_name, last_name, field_of_study, degree, specialization, mode_of_study) VALUES
        ('202311', 'Piotr', 'Wójcik', 'Informatyka', 'stopień I', 'Sztuczna Inteligencja', 'stacjonarne'),
        ('202312', 'Katarzyna', 'Pawlak', 'Telekomunikacja', 'stopień II', 'Sieci Komputerowe', 'niestacjonarne'),
        ('202313', 'Michał', 'Kaczmarek', 'Informatyka', 'stopień II', 'Bezpieczeństwo Informatyczne', 'stacjonarne'),
        ('202314', 'Agnieszka', 'Szymańska', 'Elektrotechnika', 'stopień I', 'Systemy Automatyki', 'stacjonarne'),
        ('202315', 'Bartosz', 'Kowalski', 'Inżynieria Elektroniki', 'stopień II', 'Mikroprocesory', 'niestacjonarne'),
        ('202316', 'Natalia', 'Nowicka', 'Informatyka', 'stopień I', 'Programowanie Aplikacji Mobilnych', 'stacjonarne'),
        ('202317', 'Sebastian', 'Wróbel', 'Inżynieria Elektroniki', 'stopień II', 'Komunikacja Bezprzewodowa', 'niestacjonarne'),
        ('202318', 'Martyna', 'Kwiatkowska', 'Informatyka', 'stopień I', 'Inżynieria Danych', 'stacjonarne'),
        ('202319', 'Kamil', 'Jankowski', 'Telekomunikacja', 'stopień II', 'Przetwarzanie Sygnałów', 'niestacjonarne'),
        ('202320', 'Zofia', 'Mazur', 'Informatyka', 'stopień I', 'Inżynieria Oprogramowania', 'stacjonarne'),
		('202301', 'Jan', 'Kowal', 'Informatyka', 'stopień I', 'Inżynieria Oprogramowania', 'stacjonarne'),
		('202302', 'Maria', 'Wiśniewska', 'Fizyka', 'stopień II', 'Mechanika Kwantowa', 'niestacjonarne'),
		('202303', 'Adam', 'Nowak', 'Matematyka', 'stopień II', 'Algebra', 'stacjonarne'),
		('202304', 'Ewa', 'Zielińska', 'Biologia', 'stopień I', 'Biologia Molekularna', 'stacjonarne'),
		('202305', 'Tomasz', 'Lewandowski', 'Mechanika', 'stopień II', 'Robotyka', 'niestacjonarne'),
		('202306', 'Joanna', 'Krawczyk', 'Ekonomia', 'stopień I', 'Finanse', 'stacjonarne'),
		('202307', 'Marek', 'Górski', 'Zarządzanie', 'stopień II', 'Marketing', 'niestacjonarne'),
		('202308', 'Anna', 'Lis', 'Chemia', 'stopień I', 'Chemia Analityczna', 'stacjonarne'),
		('202309', 'Karol', 'Pawlak', 'Elektrotechnika', 'stopień II', 'Energetyka', 'niestacjonarne'),
		('202310', 'Ewelina', 'Michalska', 'Inżynieria Środowiska', 'stopień I', 'Ochrona Środowiska', 'stacjonarne'),
        ('202321', 'Łukasz', 'Ostrowski', 'Informatyka', 'stopień I', 'Programowanie Gier', 'stacjonarne'),
        ('202322', 'Zuzanna', 'Pietrzyk', 'Telekomunikacja', 'stopień II', 'Systemy Mobilne', 'niestacjonarne'),
        ('202323', 'Wojciech', 'Mazurek', 'Inżynieria Elektroniki', 'stopień I', 'Elektronika Cyfrowa', 'stacjonarne'),
        ('202324', 'Aleksandra', 'Górska', 'Informatyka', 'stopień II', 'Chmura Obliczeniowa', 'stacjonarne'),
        ('202325', 'Jakub', 'Sikora', 'Telekomunikacja', 'stopień I', 'Komunikacja Satelitarna', 'niestacjonarne'),
        ('202321', 'Jakub', 'Sikora', 'Informatyka', 'stopień II', 'Sztuczna Inteligencja', 'stacjonarne'),
        ('202322', 'Marta', 'Kwiatkowska', 'Inżynieria Elektroniki', 'stopień I', 'Mikroelektronika', 'niestacjonarne'),
        ('202323', 'Tomasz', 'Nowakowski', 'Elektrotechnika', 'stopień II', 'Elektryczne Systemy Sterowania', 'stacjonarne'),
        ('202324', 'Kamil', 'Bąk', 'Telekomunikacja', 'stopień I', 'Sieci Komputerowe', 'niestacjonarne'),
        ('202325', 'Patryk', 'Zieliński', 'Informatyka', 'stopień II', 'Bezpieczeństwo Informatyczne', 'stacjonarne'),
        ('202326', 'Aleksandra', 'Sikorska', 'Inżynieria Elektroniki', 'stopień II', 'Systemy Wbudowane', 'niestacjonarne'),
        ('202327', 'Dominik', 'Majewski', 'Elektrotechnika', 'stopień I', 'Mikrosystemy', 'stacjonarne'),
        ('202328', 'Monika', 'Bąk', 'Informatyka', 'stopień II', 'Programowanie Aplikacji Webowych', 'niestacjonarne'),
        ('202329', 'Wojciech', 'Lis', 'Inżynieria Elektroniki', 'stopień I', 'Technologie Mikroprocesorowe', 'stacjonarne'),
        ('202330', 'Agnieszka', 'Chmiel', 'Informatyka', 'stopień II', 'Algorytmy i Struktury Danych', 'stacjonarne'),
        ('202331', 'Krzysztof', 'Kaczmarek', 'Inżynieria Elektroniki', 'stopień II', 'Projektowanie Układów Cyfrowych', 'niestacjonarne'),
        ('202332', 'Ewa', 'Górska', 'Telekomunikacja', 'stopień I', 'Przetwarzanie Sygnałów', 'stacjonarne'),
        ('202333', 'Szymon', 'Woźniak', 'Informatyka', 'stopień II', 'Inżynieria Danych', 'niestacjonarne'),
        ('202334', 'Zuzanna', 'Dąbrowska', 'Inżynieria Elektroniki', 'stopień I', 'Komunikacja Cyfrowa', 'stacjonarne'),
        ('202335', 'Adam', 'Mazur', 'Elektrotechnika', 'stopień II', 'Automatyka i Robotyka', 'stacjonarne'),
        ('202336', 'Julia', 'Kowalska', 'Informatyka', 'stopień I', 'Inżynieria Oprogramowania', 'niestacjonarne'),
        ('202326', 'Klaudia', 'Michałowska', 'Inżynieria Elektroniki', 'stopień II', 'Mikrosystemy Elektroniki', 'stacjonarne'),
        ('202327', 'Piotr', 'Szulc', 'Informatyka', 'stopień I', 'Rozwój Oprogramowania', 'stacjonarne'),
        ('202328', 'Katarzyna', 'Bąk', 'Telekomunikacja', 'stopień II', 'Optymalizacja Sieci', 'niestacjonarne'),
        ('202329', 'Marek', 'Walczak', 'Inżynieria Elektroniki', 'stopień I', 'Nanotechnologia', 'stacjonarne'),
        ('202327', 'Piotr', 'Kowalski', 'Informatyka', 'stopień I', 'Rozwój Oprogramowania', 'stacjonarne'),
        ('202328', 'Kamil', 'Ślimak', 'Telekomunikacja', 'stopień II', 'Optymalizacja Sieci', 'niestacjonarne'),
        ('202329', 'Mateusz', 'Gortat', 'Inżynieria Elektroniki', 'stopień I', 'Nanotechnologia', 'stacjonarne'),
        ('202330', 'Paulina', 'Kołodziej', 'Informatyka', 'stopień I', 'Technologie Internetowe', 'stacjonarne');`,

		// Insert Completed Thesis
		`INSERT INTO Completed_Thesis (thesis_number, exam_date, exam_time, average_study_grade, competency_exam_grade, diploma_exam_grade, final_study_result, final_study_result_text, thesis_title_polish, thesis_title_english, thesis_language, library, chair_id, supervisor_id, assistant_supervisor_id, reviewer_id, student_id, hourly_settlement_id, reviewer_academic_title, chair_academic_title, supervisor_academic_title, assistant_supervisor_academic_title) VALUES
		('k22/inż/001/2024', '2024-01-15', '10:00', 4.5, 5.0, 4.8, 'Very Good', 'Bardzo Dobry', 'Wpływ technologii na społeczeństwo', 'Impact of Technology on Society', 'polski', '', 8667, 8668, 8669, 8670, 2711, NULL, 'Dr.', 'Dr.', 'Mgr. inż.', 'Inż.'),
        ('k22/inż/011/2024', '2024-01-20', '10:30', 4.6, 4.7, 4.8, 'Very Good', 'Bardzo Dobry', 'Algorytmy w sztucznej inteligencji', 'Algorithms in Artificial Intelligence', 'polski', '', 8671, 8667, 8672, 8673, 2712, NULL, 'Dr.', 'Dr. hab.', 'Mgr. inż.', 'Inż.'),
        ('k22/mgr/012/2024', '2024-02-15', '13:00', 4.5, 4.6, 4.7, 'Very Good', 'Bardzo Dobry', 'Optymalizacja sieci komputerowych', 'Optimization of Computer Networks', 'angielski', '', 8674, 8671, 8676, 8675, 2713, NULL, 'Prof. dr hab.', 'Dr.', 'Dr.', 'Mgr.'),
        ('k22/inż/013/2024', '2024-03-05', '11:30', 4.8, 4.9, 5.0, 'Excellent', 'Celujący', 'Bezpieczeństwo systemów informacyjnych', 'Information Systems Security', 'polski', '', 8673, 8677, 8678, 8676, 2714, NULL, 'Dr. hab.', 'Prof. dr hab.', 'Dr.', 'Mgr. inż.'),
        ('k22/mgr/014/2024', '2024-04-10', '14:30', 4.7, 4.8, 4.9, 'Excellent', 'Celujący', 'Internet rzeczy i jego aplikacje', 'Internet of Things and its Applications', 'angielski', '', 8672, 8673, 8674, 8678, 2715, NULL, 'Dr. hab.', 'Dr.', 'Mgr.', 'Dr.'),
        ('k22/inż/015/2024', '2024-05-01', '09:00', 4.4, 4.5, 4.6, 'Good', 'Dobry', 'Mikroprocesory w elektronice', 'Microprocessors in Electronics', 'polski', '', 8676, 8671, 8677, 8673, 2716, NULL, 'Prof. dr hab.', 'Dr. hab.', 'Dr.', 'Mgr. inż.'),
        ('k22/mgr/016/2024', '2024-06-20', '12:00', 4.9, 5.0, 4.8, 'Excellent', 'Celujący', 'Przetwarzanie sygnałów w telekomunikacji', 'Signal Processing in Telecommunications', 'angielski', '', 8675, 8676, 8671, 8677, 2717, NULL, 'Dr.', 'Dr.', 'Prof. dr hab.', 'Dr. hab.'),
        ('k22/inż/017/2024', '2024-07-15', '10:00', 4.2, 4.3, 4.4, 'Good', 'Dobry', 'Zarządzanie projektami IT', 'IT Project Management', 'polski', '', 8677, 8678, 8675, 8671, 2720, NULL, 'Dr.', 'Mgr. inż.', 'Dr.', 'Inż.'),
        ('k22/mgr/018/2024', '2024-08-05', '15:30', 4.8, 4.9, 5.0, 'Excellent', 'Celujący', 'Analiza danych w informatyce', 'Data Analysis in Computer Science', 'angielski', '', 8673, 8672, 8675, 8676, 2718, NULL, 'Prof. dr hab.', 'Dr. hab.', 'Dr.', 'Mgr.'),
        ('k22/inż/019/2024', '2024-09-10', '11:30', 4.3, 4.4, 4.5, 'Good', 'Dobry', 'Systemy wbudowane w elektronice', 'Embedded Systems in Electronics', 'polski', '', 8672, 8673, 8674, 8675, 2719, NULL, 'Dr. hab.', 'Prof. dr hab.', 'Dr.', 'Mgr. inż.'),
        ('k22/inż/021/2024', '2024-01-25', '10:00', 4.6, 4.7, 4.8, 'Very Good', 'Bardzo Dobry', 'Programowanie gier komputerowych', 'Computer Game Programming', 'polski', '', 8673, 8672, 8671, 8674, 2721, NULL, 'Dr.', 'Dr. hab.', 'Mgr. inż.', 'Inż.'),
        ('k22/mgr/022/2024', '2024-02-12', '12:00', 4.7, 4.8, 4.9, 'Excellent', 'Celujący', 'Systemy mobilne w telekomunikacji', 'Mobile Systems in Telecommunications', 'angielski', '', 8674, 8673, 8672, 8671, 2722, NULL, 'Prof. dr hab.', 'Dr.', 'Dr.', 'Mgr.'),
        ('k22/inż/023/2024', '2024-03-15', '14:30', 4.9, 5.0, 5.0, 'Excellent', 'Celujący', 'Zastosowanie chmury obliczeniowej w biznesie', 'Application of Cloud Computing in Business', 'polski', '', 8672, 8671, 8676, 8673, 2723, NULL, 'Dr. hab.', 'Prof. dr hab.', 'Dr.', 'Mgr. inż.'),
        ('k22/mgr/024/2024', '2024-04-01', '13:00', 4.4, 4.5, 4.6, 'Good', 'Dobry', 'Bezpieczeństwo systemów mobilnych', 'Security of Mobile Systems', 'angielski', '', 8673, 8674, 8675, 8676, 2724, NULL, 'Dr.', 'Dr. hab.', 'Mgr.', 'Dr.'),
        ('k22/inż/025/2024', '2024-05-18', '11:00', 4.8, 4.9, 5.0, 'Excellent', 'Celujący', 'Nanotechnologia w elektronice', 'Nanotechnology in Electronics', 'polski', '', 8676, 8677, 8675, 8672, 2725, NULL, 'Prof. dr hab.', 'Dr. hab.', 'Dr.', 'Mgr. inż.'),
        ('k22/mgr/026/2024', '2024-06-25', '10:30', 4.6, 4.7, 4.8, 'Very Good', 'Bardzo Dobry', 'Optymalizacja sieci komputerowych w chmurze', 'Optimization of Computer Networks in Cloud', 'angielski', '', 8675, 8676, 8677, 8674, 2726, NULL, 'Dr.', 'Dr.', 'Prof. dr hab.', 'Dr. hab.'),
        ('k22/inż/027/2024', '2024-07-22', '15:00', 4.3, 4.4, 4.5, 'Good', 'Dobry', 'Systemy wbudowane w urządzeniach elektronicznych', 'Embedded Systems in Electronic Devices', 'polski', '', 8677, 8678, 8675, 8676, 2727, NULL, 'Dr. hab.', 'Prof. dr hab.', 'Dr.', 'Mgr. inż.'),
        ('k22/mgr/028/2024', '2024-08-05', '13:00', 4.7, 4.8, 4.9, 'Excellent', 'Celujący', 'Technologie 5G w telekomunikacji', '5G Technologies in Telecommunications', 'angielski', '', 8673, 8671, 8676, 8672, 2728, NULL, 'Prof. dr hab.', 'Dr. hab.', 'Dr.', 'Mgr.'),
        ('k22/inż/029/2024', '2024-09-15', '11:30', 4.2, 4.3, 4.4, 'Good', 'Dobry', 'Algorytmy sztucznej inteligencji w telekomunikacji', 'Artificial Intelligence Algorithms in Telecommunications', 'polski', '', 8678, 8674, 8673, 8671, 2729, NULL, 'Dr.', 'Mgr. inż.', 'Dr.', 'Inż.'),
        ('k22/mgr/030/2024', '2024-10-01', '10:00', 4.8, 4.9, 5.0, 'Excellent', 'Celujący', 'Analiza danych w systemach komputerowych', 'Data Analysis in Computer Systems', 'angielski', '', 8676, 8677, 8671, 8674, 2730, NULL, 'Dr. hab.', 'Prof. dr hab.', 'Dr.', 'Mgr.'),
        ('k22/mgr/020/2024', '2024-10-05', '13:30', 4.6, 4.7, 4.8, 'Very Good', 'Bardzo Dobry', 'Wykorzystanie algorytmów w automatyce', 'Use of Algorithms in Automation', 'angielski', '', 8674, 8671, 8676, 8673, 2731, NULL, 'Dr.', 'Dr. hab.', 'Mgr.', 'Dr.'),
		('k22/mgr/002/2024', '2024-02-20', '14:00', 4.0, 4.2, 4.3, 'Good', 'Dobry', 'Fizyka cząstek elementarnych', 'Particle Physics', 'polski', '', 8671, 8672, 8673, 8674, 2732, NULL, 'Prof. dr hab.', 'Dr. hab.', 'Dr.', 'Mgr.'),
		('k22/inż/003/2024', '2024-03-10', '12:00', 4.8, 5.0, 4.9, 'Excellent', 'Celujący', 'Zastosowania algebry', 'Applications of Algebra', 'angielski', '', 8675, 8676, 8677, 8678, 2733, NULL, 'Dr. hab.', 'Prof. dr hab.', 'Dr.', 'Mgr. inż.'),
		('k22/mgr/004/2024', '2024-04-01', '09:00', 4.6, 4.7, 4.8, 'Very Good', 'Bardzo Dobry', 'Wpływ finansów na rozwój gospodarki', 'Impact of Finance on Economic Development', 'polski', '', 8671, 8668, 8673, 8674, 2734, NULL, 'Dr.', 'Mgr. inż.', 'Dr.', 'Mgr.'),
		('k22/inż/005/2024', '2024-05-12', '11:00', 4.3, 4.4, 4.5, 'Good', 'Dobry', 'Zarządzanie ryzykiem w marketingu', 'Risk Management in Marketing', 'angielski', '', 8672, 8669, 8677, 8678, 2735, NULL, 'Prof. dr hab.', 'Dr. hab.', 'Mgr.', 'Dr.'),
		('k22/mgr/006/2024', '2024-06-18', '13:00', 4.7, 4.8, 4.9, 'Excellent', 'Celujący', 'Analiza chemiczna w badaniach środowiskowych', 'Chemical Analysis in Environmental Studies', 'polski', '', 8667, 8676, 8669, 8674, 2736, NULL, 'Dr. hab.', 'Dr.', 'Mgr. inż.', 'Inż.'),
		('k22/inż/007/2024', '2024-07-25', '10:30', 4.2, 4.3, 4.4, 'Good', 'Dobry', 'Energia odnawialna w nowoczesnej elektrotechnice', 'Renewable Energy in Modern Electrical Engineering', 'polski', '', 8675, 8668, 8677, 8673, 2737, NULL, 'Dr.', 'Dr.', 'Prof. dr hab.', 'Dr. hab.'),
		('k22/mgr/008/2024', '2024-08-15', '15:00', 4.9, 5.0, 4.8, 'Excellent', 'Celujący', 'Zastosowanie matematyki w analizie danych', 'Application of Mathematics in Data Analysis', 'angielski', '', 8671, 8667, 8672, 8678, 2738, NULL, 'Prof. dr hab.', 'Dr. hab.', 'Dr.', 'Mgr.'),
		('k22/inż/009/2024', '2024-09-05', '08:00', 4.0, 4.1, 4.2, 'Good', 'Dobry', 'Biologiczne metody oczyszczania wody', 'Biological Methods for Water Purification', 'polski', '', 8673, 8675, 8678, 8672, 2739, NULL, 'Dr.', 'Mgr. inż.', 'Dr.', 'Inż.'),
		('k22/mgr/010/2024', '2024-10-22', '16:00', 4.5, 4.6, 4.7, 'Very Good', 'Bardzo Dobry', 'Zarządzanie środowiskowe w inżynierii', 'Environmental Management in Engineering', 'angielski', '', 8676, 8671, 8667, 8675, 2740, NULL, 'Dr. hab.', 'Prof. dr hab.', 'Dr.', 'Mgr. inż.');`,

		`INSERT INTO Thesis_To_Be_Completed (thesis_number, topic_polish, topic_english, thesis_language, student_id, supervisor_academic_title, supervisor_id, assistant_supervisor_academic_title, assistant_supervisor_id) VALUES
        ('k23/inż/001/2024', 'Rozwój algorytmów w sztucznej inteligencji', 'Development of Algorithms in Artificial Intelligence', 'polski', 2740, 'Dr.', 8664, 'Mgr. inż.', 8665),
        ('k23/mgr/002/2024', 'Zastosowanie mikrosystemów w elektronice', 'Application of Microsystems in Electronics', 'polski', 2741, 'Dr. hab.', 8666, 'Dr.', 8667),
        ('k23/inż/003/2024', 'Bezpieczeństwo danych w chmurze obliczeniowej', 'Data Security in Cloud Computing', 'angielski', 2742, 'Prof. dr hab.', 8668, 'Dr.', 8669),
        ('k23/mgr/004/2024', 'Analiza danych w systemach rozproszonych', 'Data Analysis in Distributed Systems', 'angielski', 2743, 'Dr.', 8670, 'Dr. hab.', 8671),
        ('k23/inż/005/2024', 'Inżynieria oprogramowania w mikroelektronice', 'Software Engineering in Microelectronics', 'polski', 2744, 'Prof. dr hab.', 8672, 'Mgr. inż.', 8673),
        ('k23/mgr/006/2024', 'Projektowanie układów scalonych w systemach komputerowych', 'Designing Integrated Circuits in Computer Systems', 'polski', 2745, 'Dr.', 8674, 'Dr.', 8675),
        ('k23/inż/007/2024', 'Optymalizacja algorytmów w systemach wbudowanych', 'Optimization of Algorithms in Embedded Systems', 'angielski', 2746, 'Dr. hab.', 8676, 'Mgr.', 8677),
        ('k23/mgr/008/2024', 'Mikroelektronika w medycynie', 'Microelectronics in Medicine', 'polski', 2747, 'Dr.', 8665, 'Dr. hab.', 8666),
        ('k23/inż/009/2024', 'Systemy wbudowane w telekomunikacji', 'Embedded Systems in Telecommunications', 'angielski', 2748, 'Prof. dr hab.', 8671, 'Dr. hab.', 8670),
        ('k23/mgr/010/2024', 'Technologie 5G w systemach mobilnych', '5G Technologies in Mobile Systems', 'angielski', 2749, 'Dr.', 8669, 'Dr. inż.', 8673),
        ('k23/inż/011/2024', 'Analiza systemów komputerowych w edukacji', 'Analysis of Computer Systems in Education', 'polski', 2750, 'Prof. dr hab.', 8668, 'Dr.', 8674),
        ('k23/mgr/012/2024', 'Zastosowanie algorytmów uczenia maszynowego w telekomunikacji', 'Application of Machine Learning Algorithms in Telecommunications', 'angielski', 2751, 'Dr. hab.', 8672, 'Mgr. inż.', 8665),
        ('k23/inż/013/2024', 'Mikrosystemy w elektronice cyfrowej', 'Microsystems in Digital Electronics', 'polski', 2752, 'Dr.', 8666, 'Dr. hab.', 8675),
        ('k23/mgr/014/2024', 'Chmurowe systemy obliczeniowe w mikroelektronice', 'Cloud Computing Systems in Microelectronics', 'angielski', 2753, 'Prof. dr hab.', 8673, 'Dr. hab.', 8667),
        ('k23/inż/015/2024', 'Optymalizacja układów scalonych w telekomunikacji', 'Optimization of Integrated Circuits in Telecommunications', 'polski', 2754, 'Dr. inż.', 8671, 'Dr.', 8674),
        ('k23/mgr/016/2024', 'Zastosowanie algorytmów sztucznej inteligencji w systemach rozproszonych', 'Application of Artificial Intelligence Algorithms in Distributed Systems', 'angielski', 2755, 'Prof. dr hab.', 8675, 'Dr. hab.', 8669);`,
	}

	// Execute each query
	for _, query := range queries {
		slog.Info("populateDatabase", "worked", true)
		if _, err := db.Exec(query); err != nil {
			slog.Info("populateDatabase", "query", query)
			slog.Error("failed to execute query", "err", err)
			return err
		}
	}

	return nil
}
