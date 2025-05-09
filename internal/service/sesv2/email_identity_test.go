// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sesv2_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfsesv2 "github.com/hashicorp/terraform-provider-aws/internal/service/sesv2"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	itypes "github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccSESV2EmailIdentity_basic_emailAddress(t *testing.T) {
	ctx := acctest.Context(t)
	rName := acctest.RandomEmailAddress(acctest.RandomDomainName())
	resourceName := "aws_sesv2_email_identity.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SESV2ServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEmailIdentityDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEmailIdentityConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmailIdentityExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "email_identity", rName),
					acctest.MatchResourceAttrRegionalARN(ctx, resourceName, names.AttrARN, "ses", regexache.MustCompile(`identity/.+`)),
					resource.TestCheckResourceAttr(resourceName, "dkim_signing_attributes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dkim_signing_attributes.0.current_signing_key_length", ""),
					resource.TestCheckResourceAttr(resourceName, "dkim_signing_attributes.0.last_key_generation_timestamp", ""),
					resource.TestCheckResourceAttr(resourceName, "dkim_signing_attributes.0.next_signing_key_length", "RSA_1024_BIT"),
					resource.TestCheckResourceAttr(resourceName, "dkim_signing_attributes.0.signing_attributes_origin", "AWS_SES"),
					resource.TestCheckResourceAttr(resourceName, "dkim_signing_attributes.0.status", "NOT_STARTED"),
					resource.TestCheckResourceAttr(resourceName, "dkim_signing_attributes.0.tokens.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "identity_type", "EMAIL_ADDRESS"),
					resource.TestCheckResourceAttr(resourceName, "verified_for_sending_status", acctest.CtFalse),
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

func TestAccSESV2EmailIdentity_basic_domain(t *testing.T) {
	ctx := acctest.Context(t)
	rName := acctest.RandomDomainName()
	resourceName := "aws_sesv2_email_identity.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SESV2ServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEmailIdentityDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEmailIdentityConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmailIdentityExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "email_identity", rName),
					acctest.MatchResourceAttrRegionalARN(ctx, resourceName, names.AttrARN, "ses", regexache.MustCompile(`identity/.+`)),
					resource.TestCheckResourceAttr(resourceName, "dkim_signing_attributes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dkim_signing_attributes.0.current_signing_key_length", "RSA_2048_BIT"),
					acctest.CheckResourceAttrRFC3339(resourceName, "dkim_signing_attributes.0.last_key_generation_timestamp"),
					resource.TestCheckResourceAttr(resourceName, "dkim_signing_attributes.0.next_signing_key_length", "RSA_2048_BIT"),
					resource.TestCheckResourceAttr(resourceName, "dkim_signing_attributes.0.signing_attributes_origin", "AWS_SES"),
					resource.TestCheckResourceAttr(resourceName, "dkim_signing_attributes.0.status", "PENDING"),
					resource.TestCheckResourceAttr(resourceName, "dkim_signing_attributes.0.tokens.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "identity_type", "DOMAIN"),
					resource.TestCheckResourceAttr(resourceName, "verified_for_sending_status", acctest.CtFalse),
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

func TestAccSESV2EmailIdentity_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	rName := acctest.RandomEmailAddress(acctest.RandomDomainName())
	resourceName := "aws_sesv2_email_identity.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SESV2ServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEmailIdentityDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEmailIdentityConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmailIdentityExists(ctx, resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfsesv2.ResourceEmailIdentity(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccSESV2EmailIdentity_configurationSetName(t *testing.T) {
	ctx := acctest.Context(t)
	rName := acctest.RandomEmailAddress(acctest.RandomDomainName())
	resourceName := "aws_sesv2_email_identity.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SESV2ServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEmailIdentityDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEmailIdentityConfig_configurationSetName1(rName, sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmailIdentityExists(ctx, resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "configuration_set_name", "aws_sesv2_configuration_set.test1", "configuration_set_name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEmailIdentityConfig_configurationSetName2(rName, sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmailIdentityExists(ctx, resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "configuration_set_name", "aws_sesv2_configuration_set.test2", "configuration_set_name"),
				),
			},
		},
	})
}

func TestAccSESV2EmailIdentity_nextSigningKeyLength(t *testing.T) {
	ctx := acctest.Context(t)
	rName := acctest.RandomDomainName()
	resourceName := "aws_sesv2_email_identity.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SESV2ServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEmailIdentityDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEmailIdentityConfig_nextSigningKeyLength(rName, string(types.DkimSigningKeyLengthRsa2048Bit)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmailIdentityExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "dkim_signing_attributes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dkim_signing_attributes.0.next_signing_key_length", "RSA_2048_BIT"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEmailIdentityConfig_nextSigningKeyLength(rName, string(types.DkimSigningKeyLengthRsa1024Bit)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmailIdentityExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "dkim_signing_attributes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dkim_signing_attributes.0.next_signing_key_length", "RSA_1024_BIT"),
				),
			},
		},
	})
}

func TestAccSESV2EmailIdentity_domainSigning(t *testing.T) {
	ctx := acctest.Context(t)
	rName := acctest.RandomDomainName()
	resourceName := "aws_sesv2_email_identity.test"

	key1 := itypes.Base64EncodeOnce([]byte(acctest.TLSRSAPrivateKeyPEM(t, 2048)))
	selector1 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	key2 := itypes.Base64EncodeOnce([]byte(acctest.TLSRSAPrivateKeyPEM(t, 2048)))
	selector2 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SESV2ServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEmailIdentityDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEmailIdentityConfig_domainSigning(rName, key1, selector1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmailIdentityExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "dkim_signing_attributes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dkim_signing_attributes.0.domain_signing_private_key", key1),
					resource.TestCheckResourceAttr(resourceName, "dkim_signing_attributes.0.domain_signing_selector", selector1),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dkim_signing_attributes.0.domain_signing_private_key", "dkim_signing_attributes.0.domain_signing_selector"},
			},
			{
				Config: testAccEmailIdentityConfig_domainSigning(rName, key2, selector2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmailIdentityExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "dkim_signing_attributes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dkim_signing_attributes.0.domain_signing_private_key", key2),
					resource.TestCheckResourceAttr(resourceName, "dkim_signing_attributes.0.domain_signing_selector", selector2),
				),
			},
		},
	})
}

func testAccCheckEmailIdentityDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).SESV2Client(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_sesv2_email_identity" {
				continue
			}

			_, err := tfsesv2.FindEmailIdentityByID(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("SESv2 Email Identity %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckEmailIdentityExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).SESV2Client(ctx)

		_, err := tfsesv2.FindEmailIdentityByID(ctx, conn, rs.Primary.ID)

		return err
	}
}

func testAccEmailIdentityConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_sesv2_email_identity" "test" {
  email_identity = %[1]q
}
`, rName)
}

func testAccEmailIdentityConfig_configurationSetName1(rName, configurationSetName string) string {
	return fmt.Sprintf(`
resource "aws_sesv2_configuration_set" "test1" {
  configuration_set_name = %[2]q
}

resource "aws_sesv2_email_identity" "test" {
  email_identity         = %[1]q
  configuration_set_name = aws_sesv2_configuration_set.test1.configuration_set_name
}
`, rName, configurationSetName)
}

func testAccEmailIdentityConfig_configurationSetName2(rName, configurationSetName string) string {
	return fmt.Sprintf(`
resource "aws_sesv2_configuration_set" "test2" {
  configuration_set_name = %[2]q
}

resource "aws_sesv2_email_identity" "test" {
  email_identity         = %[1]q
  configuration_set_name = aws_sesv2_configuration_set.test2.configuration_set_name
}
`, rName, configurationSetName)
}

func testAccEmailIdentityConfig_nextSigningKeyLength(rName, nextSigningKeyLength string) string {
	return fmt.Sprintf(`
resource "aws_sesv2_email_identity" "test" {
  email_identity = %[1]q

  dkim_signing_attributes {
    next_signing_key_length = %[2]q
  }
}
`, rName, nextSigningKeyLength)
}

func testAccEmailIdentityConfig_domainSigning(rName, domainSigningPrivateKey, domainSigningSelector string) string {
	return fmt.Sprintf(`
resource "aws_sesv2_email_identity" "test" {
  email_identity = %[1]q

  dkim_signing_attributes {
    domain_signing_private_key = %[2]q
    domain_signing_selector    = %[3]q
  }
}
`, rName, domainSigningPrivateKey, domainSigningSelector)
}
