@datacenter
Feature: Ernest datacenter list

  Scenario: Non logged user listing
    Given I setup ernest with target "https://ernest.local"
    And I logout
    When I run ernest with "datacenter list"
    Then The output should contain "You're not allowed to perform this action, please log in"

  Scenario: Logged user listing
    Given I setup ernest with target "https://ernest.local"
    When I'm logged in as "usr" / "pwd"
    And I run ernest with "datacenter list"
    Then The output should contain "fake"
    And The output should contain "fakeaws"

  Scenario: Admin user listing
    Given I setup ernest with target "https://ernest.local"
    When I'm logged in as "ci_admin" / "pwd"
    And I run ernest with "datacenter list"
    Then The output should contain "fake"
    And The output should contain "fakeaws"

  Scenario: Listing fields
    Given I setup ernest with target "https://ernest.local"
    When I'm logged in as "usr" / "pwd"
    And I run ernest with "datacenter list"
    Then The output should contain "fake"
    And The output should contain "fakeaws"
    And The output should contain "NAME"
    And The output should contain "ID"
    And The output should contain "TYPE"
    And The output should contain "aws-fake"
    And The output should contain "vcloud-fake"


