package scaleway

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func resourceScalewayRdbUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayRdbUserCreate,
		ReadContext:   resourceScalewayRdbUserRead,
		UpdateContext: resourceScalewayRdbUserUpdate,
		DeleteContext: resourceScalewayRdbUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultRdbInstanceTimeout),
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validationUUIDorUUIDWithLocality(),
				Description:  "Instance on which the user is created",
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Database user name",
				Required:    true,
				ForceNew:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Database user password",
			},
			"is_admin": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Grant admin permissions to database user",
			},
			// Common
			"region": regionSchema(),
		},
	}
}

func resourceScalewayRdbUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rdbAPI := newRdbAPI(meta)
	// resource depends on the instance locality
	regionalID := d.Get("instance_id").(string)
	region, instanceID, err := parseRegionalID(regionalID)
	if err != nil {
		diag.FromErr(err)
	}

	ins, err := waitForRDBInstance(ctx, rdbAPI, region, instanceID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	createReq := &rdb.CreateUserRequest{
		Region:     region,
		InstanceID: ins.ID,
		Name:       d.Get("name").(string),
		Password:   d.Get("password").(string),
		IsAdmin:    d.Get("is_admin").(bool),
	}

	var user *rdb.User
	//  wrapper around StateChangeConf that will just retry write on database
	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		currentUser, errCreateUser := rdbAPI.CreateUser(createReq, scw.WithContext(ctx))
		if errCreateUser != nil {
			if is409Error(errCreateUser) {
				_, errWait := waitForRDBInstance(ctx, rdbAPI, region, ins.ID, d.Timeout(schema.TimeoutCreate))
				if errWait != nil {
					return resource.NonRetryableError(errWait)
				}
				return resource.RetryableError(errCreateUser)
			}
			return resource.NonRetryableError(errCreateUser)
		}
		// set database information
		user = currentUser
		return nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resourceScalewayRdbUserID(region, expandID(instanceID), user.Name))

	return resourceScalewayRdbUserRead(ctx, d, meta)
}

func resourceScalewayRdbUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rdbAPI := newRdbAPI(meta)
	region, instanceID, userName, err := resourceScalewayRdbUserParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = waitForRDBInstance(ctx, rdbAPI, region, instanceID, d.Timeout(schema.TimeoutRead))
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := rdbAPI.ListUsers(&rdb.ListUsersRequest{
		Region:     region,
		InstanceID: instanceID,
		Name:       &userName,
	}, scw.WithContext(ctx))

	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	var user = res.Users[0]
	_ = d.Set("instance_id", newRegionalID(region, instanceID).String())
	_ = d.Set("name", user.Name)
	_ = d.Set("is_admin", user.IsAdmin)

	d.SetId(resourceScalewayRdbUserID(region, instanceID, user.Name))

	return nil
}

func resourceScalewayRdbUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rdbAPI := newRdbAPI(meta)
	// resource depends on the instance locality
	region, instanceID, userName, err := resourceScalewayRdbUserParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = waitForRDBInstance(ctx, rdbAPI, region, instanceID, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}

	req := &rdb.UpdateUserRequest{
		Region:     region,
		InstanceID: instanceID,
		Name:       userName,
	}

	if d.HasChange("password") {
		req.Password = expandStringPtr(d.Get("password"))
	}
	if d.HasChange("is_admin") {
		req.IsAdmin = scw.BoolPtr(d.Get("is_admin").(bool))
	}

	_, err = rdbAPI.UpdateUser(req, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceScalewayRdbUserRead(ctx, d, meta)
}

func resourceScalewayRdbUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rdbAPI := newRdbAPI(meta)
	// resource depends on the instance locality
	region, instanceID, userName, err := resourceScalewayRdbUserParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = waitForRDBInstance(ctx, rdbAPI, region, instanceID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		errDeleteUser := rdbAPI.DeleteUser(&rdb.DeleteUserRequest{
			Region:     region,
			InstanceID: instanceID,
			Name:       userName,
		}, scw.WithContext(ctx))
		if errDeleteUser != nil {
			if is409Error(errDeleteUser) {
				_, errWait := waitForRDBInstance(ctx, rdbAPI, region, instanceID, d.Timeout(schema.TimeoutDelete))
				if errWait != nil {
					return resource.NonRetryableError(errWait)
				}
				return resource.RetryableError(errDeleteUser)
			}
			return resource.NonRetryableError(errDeleteUser)
		}
		// set database information
		return nil
	})

	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	return nil
}

// Build the resource identifier
// The resource identifier format is "Region/InstanceId/UserName"
func resourceScalewayRdbUserID(region scw.Region, instanceID string, userName string) (resourceID string) {
	return fmt.Sprintf("%s/%s/%s", region, instanceID, userName)
}

// Extract instance ID and username from the resource identifier.
// The resource identifier format is "Region/InstanceId/UserName"
func resourceScalewayRdbUserParseID(resourceID string) (region scw.Region, instanceID string, userName string, err error) {
	idParts := strings.Split(resourceID, "/")
	if len(idParts) != 3 {
		return "", "", "", fmt.Errorf("can't parse user resource id: %s", resourceID)
	}
	return scw.Region(idParts[0]), idParts[1], idParts[2], nil
}
