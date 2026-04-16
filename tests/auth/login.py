import os
import random
import string
import sqlite3
import requests
import pytest
import bcrypt
from datetime import datetime, UTC
from dotenv import load_dotenv


# Variaveis
load_dotenv(dotenv_path=os.path.join(os.path.dirname(__file__), "..", "..", "..", ".env"))
load_dotenv(dotenv_path=os.path.join(os.path.dirname(__file__), "..", "..", ".env"))
BASE_URL = os.getenv("BASE_URL", "http://localhost:8080")
LOGIN_ROUTE = "/ward/api/v1/auth/login"
MAX_ATTEMPTS = int(os.getenv("MAX_ATTEMPTS", "30"))
TIMEOUT = float(os.getenv("TIMEOUT", "5"))
DB_PATH = os.getenv("DB_PATH", os.path.join(os.path.dirname(__file__), "..", "..", "database.db"))


def random_string(length: int):
    chars = string.ascii_letters + string.digits
    return "".join(random.choice(chars) for _ in range(length))


def random_invalid_string(length: int):
    chars = string.ascii_letters + string.digits + string.punctuation + "     "
    return "".join(random.choice(chars) for _ in range(length))


# Criar um novo usuario
VALID_USERNAME = "test_" + random_string(10)
VALID_EMAIL = "test_" + random_string(10) + "@email.com"
VALID_PASSWORD = "Pass_" + random_string(10)

HASHED_PASSWORD = bcrypt.hashpw(VALID_PASSWORD.encode(), bcrypt.gensalt()).decode()


def insert_test_user():
    conn = sqlite3.connect(DB_PATH)
    cur = conn.cursor()

    now = datetime.now(UTC).strftime("%Y-%m-%d %H:%M:%S")

    cur.execute("""
        INSERT INTO users (username, email, password, updated_at, created_at)
        VALUES (?, ?, ?, ?, ?)
    """, (VALID_USERNAME, VALID_EMAIL, HASHED_PASSWORD, now, now))

    conn.commit()
    conn.close()


def delete_test_user():
    conn = sqlite3.connect(DB_PATH)
    cur = conn.cursor()

    cur.execute("DELETE FROM users WHERE username = ? OR email = ?", (VALID_USERNAME, VALID_EMAIL))

    conn.commit()
    conn.close()


def do_login(payload: dict):
    url = BASE_URL + LOGIN_ROUTE
    return requests.post(url, json=payload, timeout=TIMEOUT)


@pytest.fixture(scope="module", autouse=True)
def setup_and_cleanup():
    delete_test_user()
    insert_test_user()
    yield
    delete_test_user()


# Login Valido
def test_login_valid_username_password():
    payload = {
        "username": VALID_USERNAME,
        "password": VALID_PASSWORD
    }
    res = do_login(payload)
    assert res.status_code in (200, 201), res.text


def test_login_valid_email_password():
    payload = {
        "email": VALID_EMAIL,
        "password": VALID_PASSWORD
    }
    res = do_login(payload)
    assert res.status_code in (200, 201), res.text


def test_login_valid_all_fields():
    payload = {
        "username": VALID_USERNAME,
        "email": VALID_EMAIL,
        "password": VALID_PASSWORD
    }
    res = do_login(payload)
    assert res.status_code in (200, 201), res.text


# Login invalido
def test_login_missing_password():
    payload = {
        "username": VALID_USERNAME
    }
    res = do_login(payload)
    assert res.status_code in (400, 401, 422), res.text


def test_login_invalid_password():
    payload = {
        "username": VALID_USERNAME,
        "password": "senha_errada_" + random_string(8)
    }
    res = do_login(payload)
    assert res.status_code in (400, 401, 422), res.text


def test_login_invalid_email_format_no_at():
    payload = {
        "email": "testemail.com",
        "password": VALID_PASSWORD
    }
    res = do_login(payload)
    assert res.status_code in (400, 401, 422), res.text


def test_login_invalid_email_format_no_dot():
    payload = {
        "email": "test@emailcom",
        "password": VALID_PASSWORD
    }
    res = do_login(payload)
    assert res.status_code in (400, 401, 422), res.text


def test_login_username_too_long():
    payload = {
        "username": "a" * 500,
        "password": VALID_PASSWORD
    }
    res = do_login(payload)
    assert res.status_code in (400, 401, 413, 422), res.text


def test_login_password_too_long():
    payload = {
        "username": VALID_USERNAME,
        "password": "a" * 1000
    }
    res = do_login(payload)
    assert res.status_code in (400, 401, 413, 422), res.text


def test_login_random_invalid_payloads():
    for _ in range(MAX_ATTEMPTS):
        payload = {
            "username": random_invalid_string(random.randint(1, 50)),
            "email": random_invalid_string(random.randint(1, 50)),
            "password": random_invalid_string(random.randint(1, 50)),
        }

        if random.random() < 0.5:
            payload.pop(random.choice(list(payload.keys())))

        res = do_login(payload)

        assert res.status_code in (400, 401, 422), (
            f"Payload aceito errado: {payload} -> {res.status_code} {res.text}"
        )