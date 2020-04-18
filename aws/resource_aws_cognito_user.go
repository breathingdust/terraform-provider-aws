package aws

import (
	"regexp"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAwsCognitoUserPool() *schema.Resource {
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
							Type: schema.TypeString,
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
				Type: schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 128),
					validation.StringMatch(regexp.MustCompile(`^[\p{L}\p{M}\p{S}\p{N}\p{P}]+`), "must contain only non-whitespace characters"),
				),
			},
			"user_pool_id": {
				Type: schema.TypeString,
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
							Type: schema.TypeString,
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
		},
	}
}

func resourceAwsCognitoUserCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).cognitoidpconn

	// why are some properties set here, but normally they are set after calling d.GetOk()
	params := &cognitoidentityprovider.AdminCreateUserInput{
		Username: aws.String(d.Get("username").(string)),
		UserPoolId: aws.String(d.Get("user_pool_id").(string)),
	},

	if v, ok := d.GetOk("client_metadata"); ok {
		params.ClientMetadata = 
	}

	if v, ok := d.GetOk("desired_delivery_mediums"); ok {
		params.DesiredDeliveryMediums = expandStringList(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("force_alias_creation"); ok {
		params.ForceAliasCreation = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("message_action"); ok {
		params.MessageAction = aws.String(v.(string)),
	}
	
	if v, ok := d.GetOk("temporary_password"); ok {
		params.MessageAction = aws.String(v.(string)),
	}	

	// how to do complex 
}

func resourceAwsCognitoUserRead(d *schema.ResourceData, meta interface{}) error {
}

func resourceAwsCognitoUserUpdate(d *schema.ResourceData, meta interface{}) error {
}

func resourceAwsCognitoUserDelete(d *schema.ResourceData, meta interface{}) error {
}
