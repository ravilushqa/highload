Feature: Testing Social Network Flow

  Scenario: Main Page
    Given I make a GET request to "/"
    Then the response code should be 200

  Scenario: Register First User
    Given I make a POST request to "/register" with the following form data:
      | register-form-first-name | register-form-last-name | register-form-password | register-form-email | register-form-birthday | register-form-sex |
      | Ravil                    | Galaktionov             | pass1                  | user1@test.io       | 2000-01-01             | 1                 |
    Then the response code should be 200

  Scenario: Register Second User
    Given I make a POST request to "/register" with the following form data:
      | register-form-first-name | register-form-last-name | register-form-password | register-form-email | register-form-birthday | register-form-sex |
      | John                     | Doe                     | pass2                  | user2@test.io       | 2010-01-01             | 2                 |
    Then the response code should be 200

  Scenario: Login First User
    Given I make a POST request to "/login" with the following form data:
      | login-form-email | login-form-password |
      | user1@test.io    | pass1               |
    Then the response code should be 200

  Scenario: Login Bad Credentials
    Given I make a POST request to "/login" with the following form data:
      | login-form-email | login-form-password |
      | user1@test.io    | bad                 |
    Then the response code should be 401

