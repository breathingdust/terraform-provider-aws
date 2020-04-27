package aws

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	// resource.AddTestSweepers("aws_cognito_user", &resource.Sweeper{
	// 	Name: "aws_cognito_user",
	// 	F:    testSweepCognitoUserPools,
	// })
}

func TestAccAWSCognitoUser_basic(t *testing.T) {
	resourceName := "aws_cognito_user.test"
	poolName := fmt.Sprintf("tf-acc-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	username := fmt.Sprintf("tf-acc-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSCognitoIdentityProvider(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSCognitoUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSCognitoUser_UserName(poolName, username),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSCognitoUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttrSet(resourceName, "user_create_date"),
					resource.TestCheckResourceAttrSet(resourceName, "user_last_modified_date"),
					resource.TestCheckResourceAttrSet(resourceName, "user_status"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAWSCognitoUserExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id := rs.Primary.Attributes["username"]
		userPoolID := rs.Primary.Attributes["user_pool_id"]

		if id == "" {
			return errors.New("No Cognito Username set")
		}

		if userPoolID == "" {
			return errors.New("No Cognito User Pool Id set")
		}

		conn := testAccProvider.Meta().(*AWSClient).cognitoidpconn

		params := &cognitoidentityprovider.AdminGetUserInput{
			Username:   aws.String(rs.Primary.Attributes["username"]),
			UserPoolId: aws.String(rs.Primary.Attributes["user_pool_id"]),
		}

		_, err := conn.AdminGetUser(params)
		return err
	}
}

func testAccCheckAWSCognitoUserDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).cognitoidpconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_cognito_user" {
			continue
		}

		params := &cognitoidentityprovider.AdminGetUserInput{
			Username:   aws.String(rs.Primary.ID),
			UserPoolId: aws.String(rs.Primary.Attributes["user_pool_id"]),
		}

		_, err := conn.AdminGetUser(params)

		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == cognitoidentityprovider.ErrCodeResourceNotFoundException {
				return nil
			}
			return err
		}
	}

	return nil
}

func testAccAWSCognitoUser_UserName(poolName, username string) string {
	return fmt.Sprintf(`
resource "aws_cognito_user_pool" "main" {
	name = "%s"
}

resource "aws_cognito_user" "test" {
  username = "%s"
  user_pool_id = "${aws_cognito_user_pool.main.id}"
}
`, poolName, username)
}
