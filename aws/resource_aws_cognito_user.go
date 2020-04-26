package aws

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAwsCognitoUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsCognitoUserCreate,
		Read:   resourceAwsCognitoUserRead,
		Update: resourceAwsCognitoUserUpdate,
		Delete: resourceAwsCognitoUserDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminCreateUser.html
		Schema: map[string]*schema.Schema{
			"client_metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"desired_delivery_mediums": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						cognitoidentityprovider.DeliveryMediumTypeSms,
						cognitoidentityprovider.DeliveryMediumTypeEmail,
					}, false),
				},
			},
			"force_alias_creation": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"message_action": {
				Type:     schema.TypeString,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						cognitoidentityprovider.MessageActionTypeResend,
						cognitoidentityprovider.MessageActionTypeSuppress,
					}, false),
				},
			},
			"temporary_password": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(6, 256),
					validation.StringMatch(regexp.MustCompile(`^[[\S]+`), "must contain only non-whitespace characters"),
				),
			},
			"user_attributes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.All(
								validation.StringLenBetween(1, 32),
								// looks like the go sdk has validators in validators.go
								validation.StringMatch(regexp.MustCompile(`^[\p{L}\p{M}\p{S}\p{N}\p{P}]+`), "must contain only non-whitespace characters"),
							),
						},
						"value": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 2048),
						},
					},
				},
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 128),
					validation.StringMatch(regexp.MustCompile(`^[\p{L}\p{M}\p{S}\p{N}\p{P}]+`), "must contain only non-whitespace characters"),
				),
			},
			"user_pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 55),
					validation.StringMatch(regexp.MustCompile(`^[\w-]+_[0-9a-zA-Z]`), "must contain alphanumeric characters, dashes or underscores"),
				),
			},
			"validation_data": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.All(
								validation.StringLenBetween(1, 32),
								validation.StringMatch(regexp.MustCompile(`^[\p{L}\p{M}\p{S}\p{N}\p{P}]+`), "must contain only non-whitespace characters"),
							),
						},
						"value": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 2048),
						},
					},
				},
			},
			"user_create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_last_modified_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAwsCognitoUserCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).cognitoidpconn

	// why are some properties set with Get() rather than GetOk(). I understand GetOk() allows you to check if there is a value.
	// Is it to do with option/required fields? Is it to allow default values for optional?
	// we do we have both Optional and Required attributes in the schema?
	params := &cognitoidentityprovider.AdminCreateUserInput{
		Username:   aws.String(d.Get("username").(string)),
		UserPoolId: aws.String(d.Get("user_pool_id").(string)),
	}

	if v, ok := d.GetOk("client_metadata"); ok {
		params.ClientMetadata = stringMapToPointers(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("desired_delivery_mediums"); ok {
		params.DesiredDeliveryMediums = expandStringList(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("force_alias_creation"); ok {
		params.ForceAliasCreation = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("message_action"); ok {
		params.MessageAction = aws.String(v.(string))
	}

	if v, ok := d.GetOk("temporary_password"); ok {
		params.TemporaryPassword = aws.String(v.(string))
	}

	if v, ok := d.GetOk("user_attributes"); ok {
		ua, err := expandAttributeTypes(v.([]interface{}))
		if err != nil {
			return err
		}
		params.UserAttributes = ua
	}

	if v, ok := d.GetOk("validation_data"); ok {
		ua, err := expandAttributeTypes(v.([]interface{}))
		if err != nil {
			return err
		}
		params.ValidationData = ua
	}

	log.Printf("[DEBUG] Creating Cognito User: %s", params)

	resp, err := conn.AdminCreateUser(params)

	if err != nil {
		return fmt.Errorf("Error creating Cognito User: %s", err)
	}

	d.SetId(*resp.User.Username)

	return resourceAwsCognitoUserGroupRead(d, meta)
}

func resourceAwsCognitoUserRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).cognitoidpconn

	params := &cognitoidentityprovider.AdminGetUserInput{
		Username:   aws.String(d.Id()),
		UserPoolId: aws.String(d.Get("user_pool_id").(string)),
	}

	resp, err := conn.AdminGetUser(params)

	// There is also a UserNotFoundException. Which to use?
	if isAWSErr(err, cognitoidentityprovider.ErrCodeResourceNotFoundException, "") {
		log.Printf("[WARN] Cognito User (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error describing Cognito User (%s): %w", d.Id(), err)
	}

	// How reconcile properties that exist on the create but not the get.?

	d.Set("enabled", resp.Enabled)

	if err := d.Set("user_create_date", resp.UserCreateDate.Format(time.RFC3339)); err != nil {
		return err
	}

	if err := d.Set("user_last_modified_date", resp.UserLastModifiedDate.Format(time.RFC3339)); err != nil {
		return err
	}
	return err
}

func resourceAwsCognitoUserUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAwsCognitoUserDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func expandAttributeTypes(s []interface{}) ([]*cognitoidentityprovider.AttributeType, error) {
	if len(s) == 0 {
		return nil, nil
	}
	attributeTypes := make([]*cognitoidentityprovider.AttributeType, 0)
	for _, raw := range s {
		p := raw.(map[string]interface{})
		n := p["name"].(string)
		v := p["value"].(string)

		attributeType := &cognitoidentityprovider.AttributeType{
			Name:  aws.String(n),
			Value: aws.String(v),
		}
		attributeTypes = append(attributeTypes, attributeType)
	}
	return attributeTypes, nil
}
