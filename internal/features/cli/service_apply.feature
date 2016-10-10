@service @service_apply
Feature: Service apply

  Scenario: Non logged service apply
    Given I setup ernest with target "https://ernest.local"
    And I logout
    When I run ernest with "service apply"
    Then The output should contain "You're not allowed to perform this action, please log in"
    When I run ernest with "service apply definitions/aws1.yml"
    Then The output should contain "You're not allowed to perform this action, please log in"

  Scenario: Logged service apply errors
    Given I setup ernest with target "https://ernest.local"
    And I'm logged in as "usr" / "pwd"
    When I run ernest with "service apply"
    Then The output should contain "You should specify a valid template path or store an ernest.yml on the current folder"
    When I run ernest with "service apply internal/definitions/unexisting_dc.yml"
    Then The output should contain "Specified datacenter does not exist"

  Scenario: Logged service apply
    Given I setup ernest with target "https://ernest.local"
    And I'm logged in as "usr" / "pwd"
    And the datacenter "test_dc" does not exist
    And I run ernest with "datacenter create aws --token tmp_token --secret tmp_secret --region tmp_region --fake test_dc"
    And the service "aws_test_service" does not exist
    When I run ernest with "service apply internal/definitions/aws1.yml"
    Then The output line number "4" should contain "VPC created"
    And The output line number "5" should contain "Creating networks:"
    And The output line number "6" should contain "test_dc-aws_test_service-web"
    And The output line number "7" should contain "IP     : 10.1.0.0/24"
    And The output line number "8" should contain "AWS ID : foo"
    And The output line number "9" should contain "Status : completed"
    And The output line number "10" should contain "Networks successfully created"
    And The output line number "11" should contain "Setting up firewalls:"
    And The output line number "12" should contain "Firewalls Created"
    And The output line number "13" should contain "Creating instances:"

