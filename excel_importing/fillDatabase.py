import pandas as pd
import sqlite3
import re

PATH_TO_EXCEL = 'Path/To/baza mock egzaminow dypl.xlsx'
PATH_TO_DATABASE = 'Path/To/diploma_database.db'

specializations = []
fields_of_study = []


def parse_names(df, name, index):
    full_name = df[name].iloc[index]

    if pd.notna(full_name):
        name_parts = full_name.split(' ', 1)
        if len(name_parts) == 2:
            return name_parts[0], name_parts[1]
        else:
            return name_parts[0], None
    else:
        return None, None


def standardise_degree(degree_str):
    if pd.notna(degree_str):
        degrees = {
            'inżynierskiego': 'I stopień',
            'magisterskiego': 'II stopień',
            'magisterskiego ': 'II stopień',
            'magisterskiego uzupełniającego': 'II stopień',
            'inżynierskiego ': 'I stopień',
            'mgr uzup.': 'II stopień',
            'magisterskiego uzup.': 'II stopień',
            'mgr': 'II stopień',
            'inż.': 'I stopień',
            'inż..': 'I stopień',
            'inż. 1o': 'I stopień',
            'mgr ': 'II stopień',
            'mg': 'II stopień',
            'inż. ': 'I stopień',
            '-': None, 'mgr.': 'II stopień'
        }

        degree_str = degree_str.lower()

        degree = degrees.get(degree_str)

        return degree
    else:
        return None


def standardize_title(title_str):
    if pd.notna(title_str):
        titles = {
            'dr hab inż': 'dr hab. inż.',
            'dr inż': 'dr inż.',
            'dr hab inz - prof pł': 'dr hab. inż.',
            'prof dr hab inż': 'Prof. dr hab. inż.',
            'dr hab inż - prof pł': 'dr hab. inż.',
            'mgr inż': 'mgr inż.',
            'dr inż + mgr inż': 'dr inż. + mgr inż.',  # ?
            'prof dr hab inż + mgr': 'Prof. dr hab. inż. + mgr',  # ?
            'dr inz': 'dr inż.',
            'x': None,
            'prof dr hab inż i mgr inż': 'Prof. dr hab. inż. + mgr inż.',  # ?
            'dr hab inz': 'dr hab. inż.',
            'bez udz dypl': None,
            'dr inż szczepaniak jakub': 'dr inż.',
            'dr': 'dr',
            'dr hab': 'dr hab.',
            'mg inż': 'mgr inż.',
            'prof': 'Prof.',
            '-': None,
            'dr hab - prof pł': 'dr hab.',
            'd inż': 'dr inż.',
            'mgr': 'mgr',
            '- - -': None,
            'prof dr hab': 'Prof. dr hab.'
        }
        title_str = title_str.lower().replace(' -', ' - ') \
            .replace('-', ' - ') \
            .replace('+m', '+ m') \
            .replace('.', ' ') \
            .replace('   ', ' ') \
            .replace('  ', ' ').strip()

        title = titles.get(title_str)

        return title
    else:
        return None


def parse_exam_date(date_str):
    if pd.notna(date_str):
        months = {
            "stycznia": "01", "styczeń": "01", "lutego": "02", "luty": "02", "marca": "03",
            "marzec": "03", "kwietnia": "04", "kwietna": "04", "kwetnia": "04", "kwiecień": "04",
            "maja": "05", "maj": "05", "czerwca": "06", "czerwiec": "06", "lipca": "07",
            "lipiec": "07", "sierpnia": "08", "sierpień": "08", "września": "09", "wrzesnia": "09",
            "wrz": "09", "wrzesień": "09", "października": "10", "październik": "10", "paźdz": "10",
            "listopada": "11", "listopad": "11", "list": "11", "grudnia": "12", "grudzień": "12", "grud": "12"
        }

        date_str = date_str.replace('.', ' ')

        match = re.search(r"(\d+)\s*([a-zA-Zśź]+)\s*(\d{4})", date_str)
        if match:
            day = match.group(1)
            month_str = match.group(2).strip().lower()
            year = match.group(3)

            month = months.get(month_str)

            if month:
                return f"{year}-{month.zfill(2)}-{day.zfill(2)}"
            else:
                print("Incorect date:", date_str)
                return None
        else:
            return None
    else:
        return None


def create_comment(place_of_birth, date_of_birth, son_daughter, immatriculation_year):
    comment = ''
    if pd.notna(place_of_birth):
        comment += 'Miejsce urodzenia: ' + str(place_of_birth) + '; '
    if pd.notna(date_of_birth):
        comment += 'Data urodzenia: ' + str(date_of_birth) + '; '
    if pd.notna(son_daughter):
        comment += 'Syn/córka: ' + str(son_daughter) + '; '
    if pd.notna(immatriculation_year):
        comment += 'Rok immatrykulacji: ' + str(int(immatriculation_year))
    if pd.isna(place_of_birth) and pd.isna(date_of_birth) and pd.isna(son_daughter) and pd.isna(immatriculation_year):
        comment = None

    return comment


def standardize_study_mode(mode_of_study_str):
    if pd.notna(mode_of_study_str):
        modes = {
            'WYPOŻ.': None, 'stac.': 'stacjonarne', 'niestac.': 'niestacjonarne',
            'stac': 'stacjonarne', 'niest.': 'niestacjonarne'
        }
        mode_of_study_str = mode_of_study_str.strip()
        mode_of_study = modes.get(mode_of_study_str)

        return mode_of_study
    else:
        return None


def get_or_create_student(df, index, cursor):
    print(df[['nazwisko', 'syn_córka', 'ur_dnia', 'miejsce', 'nr_albumu', 'kierunek', 'specjalnosc', 'Tryb stud.',
              'MgrInż', 'rok_immatrykulacji']].iloc[index])
    print()

    first_name, last_name = parse_names(df, 'nazwisko', index)

    if not first_name and not last_name:
        return None

    cursor.execute('''
                SELECT id FROM Student 
                WHERE student_number = ? AND first_name = ? AND last_name = ?
            ''', (df['nr_albumu'].iloc[index], first_name, last_name))

    student = cursor.fetchone()

    if student:
        student_id = student[0]
    else:
        comment = create_comment(df['miejsce'].iloc[index], df['ur_dnia'].iloc[index],
                                 df['syn_córka'].iloc[index], df['rok_immatrykulacji'].iloc[index])

        mode_of_study = standardize_study_mode(df['Tryb stud.'].iloc[index])

        degree = standardise_degree(df['MgrInż'].iloc[index])

        if pd.notna(df['nr_albumu'].iloc[index]):
            student_number = str(df['nr_albumu'].iloc[index]).strip()
        else:
            student_number = None

        if pd.notna(df['kierunek'].iloc[index]):

            field_of_study = str(df['kierunek'].iloc[index]).strip()
            if field_of_study not in fields_of_study:
                fields_of_study.append(field_of_study)
        else:
            field_of_study = None

        if pd.notna(df['specjalnosc'].iloc[index]):
            specialization = str(df['specjalnosc'].iloc[index]).strip()
            if specialization not in specializations:
                specializations.append(specialization)
        else:
            specialization = None

        cursor.execute('''
                        INSERT INTO Student (student_number, first_name, last_name, field_of_study,
                        specialization, mode_of_study, comment, degree)
                        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
                    ''', (
            student_number,
            first_name,
            last_name,
            field_of_study,
            specialization,
            mode_of_study,
            comment,
            degree
        ))
        student_id = cursor.lastrowid

    return student_id


def get_or_create_employee(df, index, name, title, cursor):
    print(df[[title, name]].iloc[index])
    print()

    first_name, last_name = parse_names(df, name, index)

    if not first_name and not last_name:
        return None

    academic_title = standardize_title(df[title].iloc[index])

    cursor.execute('''
                SELECT id FROM University_Employee 
                WHERE first_name = ? AND last_name = ? AND current_academic_title = ?
            ''', (first_name, last_name, academic_title))

    employee = cursor.fetchone()

    if employee:
        employee_id = employee[0]
    else:
        cursor.execute('''
                        INSERT INTO University_Employee (current_academic_title, first_name, last_name)
                        VALUES (?, ?, ?)
                    ''', (
            academic_title,
            first_name,
            last_name
        ))

        employee_id = cursor.lastrowid

    return employee_id


def fill_specializations_and_fields_of_study(cursor):
    for specialization in specializations:
        cursor.execute('''
                INSERT INTO Specializations (name)
                VALUES (?)
            ''', (
            specialization,
        ))

    for field_of_study in fields_of_study:
        cursor.execute('''
                INSERT INTO fields_of_study (name)
                VALUES (?)
            ''', (
            field_of_study,
        ))


def main():
    df = pd.read_excel(PATH_TO_EXCEL, engine='openpyxl')

    conn = sqlite3.connect(PATH_TO_DATABASE)
    cursor = conn.cursor()

    print("Kolumny dostępne w Excelu:", df.columns)

    for index, row in df.iterrows():
        if row.isnull().all():
            print("Skipped row:")
            print(row)
            continue

        if any(str(cell).strip().isdigit() and len(str(cell).strip()) == 4 for cell in row):
            print("Skipped row:")
            print(row)
            continue

        if row['średnia ze studiów'] == 'średnia z ocen':
            print("Skipped row:")
            print(row)
            continue

        # print('student')
        student_id = get_or_create_student(df, index, cursor)

        # print('chair')
        chair_id = get_or_create_employee(df, index, 'przewodniczący', 'tytuł', cursor)

        # print('supervisor')
        supervisor_id = get_or_create_employee(df, index, 'Promotor', 'tytuł.1', cursor)

        # print('assistant supervisor')
        assistant_supervisor_id = get_or_create_employee(df, index, 'opiekun', 'Unnamed: 18', cursor)

        # print('reviewer')
        reviewer_id = get_or_create_employee(df, index, 'Recenzent', 'tytuł.2', cursor)

        print('praca')
        print(df[['Nr pracy', 'data egz.', 'średnia ze studiów', 'Ocena pracydypl.',
                  'Temat pracy dyplomowej', 'Temat pracy dyplomowej w języku angielskim']].iloc[index])

        if pd.notna(df['Nr pracy'].iloc[index]):
            thesis_number = str(df['Nr pracy'].iloc[index]).strip()
        else:
            thesis_number = None

        if pd.notna(df['średnia ze studiów'].iloc[index]):
            average_from_study = str(df['średnia ze studiów'].iloc[index]).strip()
        else:
            average_from_study = None

        if pd.notna(df['Ocena pracydypl.'].iloc[index]):
            thesis_grade = str(df['Ocena pracydypl.'].iloc[index]).strip()
        else:
            thesis_grade = None

        if pd.notna(df['Temat pracy dyplomowej'].iloc[index]):
            topic_pl = str(df['Temat pracy dyplomowej'].iloc[index]).strip()
        else:
            topic_pl = None

        if pd.notna(df['Temat pracy dyplomowej w języku angielskim'].iloc[index]):
            topic_en = str(df['Temat pracy dyplomowej w języku angielskim'].iloc[index]).strip()
        else:
            topic_en = None

        chair_academic_title = standardize_title(df['tytuł'].iloc[index])
        supervisor_academic_title = standardize_title(df['tytuł.1'].iloc[index])
        assistant_supervisor_academic_title = standardize_title(df['Unnamed: 18'].iloc[index])
        reviewer_academic_title = standardize_title(df['tytuł.2'].iloc[index])

        cursor.execute('''
                INSERT INTO Completed_Thesis (thesis_number, exam_date, average_study_grade,
                final_study_result, thesis_title_polish, thesis_title_english, student_id,
                chair_id, supervisor_id, assistant_supervisor_id, reviewer_id, chair_academic_title,
                supervisor_academic_title, assistant_supervisor_academic_title, reviewer_academic_title)
                VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
            ''', (
            thesis_number,
            parse_exam_date(df['data egz.'].iloc[index]),
            average_from_study,
            thesis_grade,
            topic_pl,
            topic_en,
            student_id,
            chair_id,
            supervisor_id,
            assistant_supervisor_id,
            reviewer_id,
            chair_academic_title,
            supervisor_academic_title,
            assistant_supervisor_academic_title,
            reviewer_academic_title
        ))

    fill_specializations_and_fields_of_study(cursor)

    conn.commit()
    conn.close()


if __name__ == '__main__':
    main()
