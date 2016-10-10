@service
Feature: Ernest service

  Scenario: running "service" command as logged in user
    Given I setup ernest with target "https://ernest.local"
    When I'm logged in as "usr" / "pwd"
    And I run ernest with "service"
    And The output should contain "list"
    And The output should contain "create"

  Scenario: running "service" command as non logged in user
    Given I setup ernest with target "https://ernest.local"
    And I logout
    When I run ernest with "group"
    And The output should contain "list"
    And The output should contain "create"

