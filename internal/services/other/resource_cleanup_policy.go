package other

import (
	"encoding/json"
	"fmt"

	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	nexusSchema "github.com/datadrivers/go-nexus-client/nexus3/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/serialt/terraform-provider-nexus/internal/schema/common"
)

func ResourceCleanUpPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a Nexus Cleanup Policy Rule.",

		Create: resourceCleanUpPolicyCreate,
		Read:   resourceCleanUpPolicyRead,
		Update: resourceCleanUpPolicyUpdate,
		Delete: resourceCleanUpPolicyDelete,
		Exists: resourceCleanUpPolicyExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": common.ResourceID,
			"name": {
				Description: "The name of the cleanup policy rule",
				ForceNew:    true,
				Type:        schema.TypeString,
				Required:    true,
			},
			"format": {
				Description: "The format that this cleanup policy can be applied to",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},
			"notes": {
				Description: "Notes for this policy",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"criteria": {
				Description: "Cleanup criteria",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_downloaded_days": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Remove components that were published over this amount of time",
						},
						"last_blob_updated_days": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Remove components that haven't been downloaded in this amount of time",
						},
						"regex": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "Remove components that have at least one asset name matching the following" +
								" regular expression pattern",
						},
					},
				},
			},
		},
	}
}

type CleanUpPolicy struct {
	Name     string   `json:"name"`
	Format   string   `json:"format"`
	Notes    string   `json:"notes"`
	Criteria Criteria `json:"criteria"`
}

type Criteria struct {
	LastDownloaded  int    `json:"lastDownloaded"`
	LastBlobUpdated int    `json:"lastBlobUpdated"`
	Regex           string `json:"regex"`
}

func getCleanUpPolicyFromResourceData(d *schema.ResourceData) CleanUpPolicy {
	c := CleanUpPolicy{
		Name:   d.Get("name").(string),
		Format: d.Get("format").(string),
		Notes:  d.Get("notes").(string),
	}
	if _, ok := d.GetOk("criteria"); ok {
		rList := d.Get("criteria").([]interface{})
		rConfig := rList[0].(map[string]interface{})
		c.Criteria = Criteria{
			LastDownloaded:  rConfig["last_downloaded_days"].(int),
			LastBlobUpdated: rConfig["last_blob_updated_days"].(int),
			Regex:           rConfig["regex"].(string),
		}
	}
	return c
}

func GetCreationPayLoad(cu *CleanUpPolicy) string {
	r, _ := json.Marshal(cu)
	return string(r)
}

func NewCleanUpScript(name string) nexusSchema.Script {
	var content = `// Original from:
// https://github.com/idealista/nexus-role/blob/master/files/scripts/cleanup_policy.groovy
import com.google.common.collect.Maps
import groovy.json.JsonSlurper
import groovy.json.JsonBuilder
import java.util.concurrent.TimeUnit

import org.sonatype.nexus.cleanup.storage.CleanupPolicyStorage
import static org.sonatype.nexus.cleanup.config.CleanupPolicyConstants.LAST_BLOB_UPDATED_KEY
import static org.sonatype.nexus.cleanup.config.CleanupPolicyConstants.LAST_DOWNLOADED_KEY
import static org.sonatype.nexus.cleanup.config.CleanupPolicyConstants.REGEX_KEY


def cleanupPolicyStorage = container.lookup(CleanupPolicyStorage.class.getName())

try {
    parsed_args = new JsonSlurper().parseText(args)
} catch(Exception ex) {
    log.debug("list")
    def policies = []
    cleanupPolicyStorage.getAll().each {
        policies << toJsonString(it)
    }
    return policies
}

parsed_args.each {
    log.debug("Received arguments: <${it.key}=${it.value}> (${it.value.getClass()})")
}

if (parsed_args.name == null) {
    throw new Exception("Missing mandatory argument: name")
}

// "get" operation
if (parsed_args.size() == 1) {
    log.debug("get")
    existingPolicy = cleanupPolicyStorage.get(parsed_args.name)
    return toJsonString(existingPolicy)
}

// create and update use this
Map<String, String> criteriaMap = createCriteria(parsed_args)

// "update" operation
if (cleanupPolicyStorage.exists(parsed_args.name)) {
    log.debug("Updating Cleanup Policy <name=${parsed_args.name}>")
    existingPolicy = cleanupPolicyStorage.get(parsed_args.name)
    existingPolicy.setNotes(parsed_args.notes)
    existingPolicy.setCriteria(criteriaMap)
    cleanupPolicyStorage.update(existingPolicy)
    return toJsonString(existingPolicy)
}

// "create" operation
format = parsed_args.format == "all" ? "ALL_FORMATS" : parsed_args.format

log.debug("Creating Cleanup Policy <name=${parsed_args.name}>")
cleanupPolicy = cleanupPolicyStorage.newCleanupPolicy()

log.debug("Configuring Cleanup Policy <policy=${cleanupPolicy}>")
cleanupPolicy.setName(parsed_args.name)
cleanupPolicy.setNotes(parsed_args.notes)
cleanupPolicy.setFormat(format)
cleanupPolicy.setMode('delete')
cleanupPolicy.setCriteria(criteriaMap)

log.debug("Adding Cleanup Policy <policy=${cleanupPolicy}>")
cleanupPolicyStorage.add(cleanupPolicy)
return toJsonString(cleanupPolicy)


def Map<String, String> createCriteria(parsed_args) {
    Map<String, String> criteriaMap = Maps.newHashMap()
    if (parsed_args.criteria.lastBlobUpdated == null) {
        criteriaMap.remove(LAST_BLOB_UPDATED_KEY)
    } else {
        criteriaMap.put(LAST_BLOB_UPDATED_KEY, asStringSeconds(parsed_args.criteria.lastBlobUpdated))
    }
    if (parsed_args.criteria.lastDownloaded == null) {
        criteriaMap.remove(LAST_DOWNLOADED_KEY)
    } else {
        criteriaMap.put(LAST_DOWNLOADED_KEY, asStringSeconds(parsed_args.criteria.lastDownloaded))
    }
    if (parsed_args.criteria.regex == null) {
        criteriaMap.remove(REGEX_KEY)
    } else {
        criteriaMap.put(REGEX_KEY, parsed_args.criteria.regex)
    }
    log.debug("Using criteriaMap: ${criteriaMap}")

    return criteriaMap
}

def Integer asSeconds(days) {
    return days * TimeUnit.DAYS.toSeconds(1)
}

def String asStringSeconds(daysInt) {
    return String.valueOf(asSeconds(daysInt))
}

// There's got to be a better way to do this.
// using JsonOutput directly on the object causes a stack overflow
def String toJsonString(cleanup_policy) {
    def policyString = new JsonBuilder()
    policyString {
        name cleanup_policy.getName()
        notes cleanup_policy.getNotes()
        format cleanup_policy.getFormat()
        mode cleanup_policy.getMode()
        criteria cleanup_policy.getCriteria()
    }
    return policyString.toPrettyString()
}
`
	return nexusSchema.Script{
		Name:    name,
		Content: content,
		Type:    "groovy",
	}
}

func NewDeleteCleanUpScript(name string) nexusSchema.Script {
	content := `
import groovy.json.JsonSlurper
import org.sonatype.nexus.cleanup.storage.CleanupPolicyStorage
import org.sonatype.nexus.cleanup.storage.CleanupPolicyComponent


parsed_args = new JsonSlurper().parseText(args)
if (parsed_args.name == null) {
    throw new Exception("Missing mandatory argument: name")
}
def deleteCleanupPolicy(String name) {
    def cleanupPolicyStorage = container.lookup(CleanupPolicyStorage.class.getName())
    def cleanupPolicyComponent = container.lookup(CleanupPolicyComponent.class.getName())
    if (cleanupPolicyStorage.exists(name)) {
        cleanupPolicyStorage.remove(cleanupPolicyStorage.get(name))
    }
}
deleteCleanupPolicy(parsed_args.name)
`
	return nexusSchema.Script{
		Name:    fmt.Sprintf("%s-delete", name),
		Content: content,
		Type:    "groovy",
	}
}
func resourceCleanUpPolicyCreate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	cu := getCleanUpPolicyFromResourceData(resourceData)
	createScript := NewCleanUpScript(cu.Name)
	//creation script
	if err := client.Script.Create(&createScript); err != nil {
		return err
	}
	//run creation script
	payload := GetCreationPayLoad(&cu)
	fmt.Printf("payload: %v", payload)
	if err := client.Script.RunWithPayload(createScript.Name, payload); err != nil {
		return err
	}

	//create deletion script
	deletionScript := NewDeleteCleanUpScript(cu.Name)
	if err := client.Script.Create(&deletionScript); err != nil {
		return err
	}

	//cleanup policy name is equal to script name
	resourceData.SetId(cu.Name)

	return resourceCleanUpPolicyRead(resourceData, m)
}

func resourceCleanUpPolicyRead(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	script, err := client.Script.Get(resourceData.Id())
	if err != nil {
		return err
	}

	if script == nil {
		resourceData.SetId("")
		return nil
	}
	/**Instead of set data from what we get from nexus
	We set data to what are provided in resource data
	**/
	resourceData.Set("name", script.Name)
	cu := getCleanUpPolicyFromResourceData(resourceData)
	resourceData.Set("format", cu.Format)
	resourceData.Set("notes", cu.Notes)
	resourceData.Set("criteria", flattenCriteria(&cu.Criteria))
	return nil
}

func resourceCleanUpPolicyUpdate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	// if name is changed, resource is force to recreate already
	if resourceData.HasChangeExcept("name") {
		cu := getCleanUpPolicyFromResourceData(resourceData)
		script := NewCleanUpScript(cu.Name)
		if err := client.Script.Update(&script); err != nil {
			return err
		}
		payload := GetCreationPayLoad(&cu)
		if err := client.Script.RunWithPayload(cu.Name, payload); err != nil {
			return err
		}
	}
	return resourceCleanUpPolicyRead(resourceData, m)
}

// there is no api for cleanup policy, so we delete it with also a script
func resourceCleanUpPolicyDelete(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	//Delete creation script
	if err := client.Script.Delete(resourceData.Id()); err != nil {
		return err
	}
	//Run delete cleanup policy script
	cu := getCleanUpPolicyFromResourceData(resourceData)
	deletionScript := NewDeleteCleanUpScript(cu.Name)
	deletionPayload := fmt.Sprintf(`
	{
		"name": "%s",
	}
`, cu.Name)
	if err := client.Script.RunWithPayload(deletionScript.Name, deletionPayload); err != nil {
		return err
	}
	//Delete deletion script
	if err := client.Script.Delete(deletionScript.Name); err != nil {
		return err
	}
	resourceData.SetId("")
	return nil
}

func resourceCleanUpPolicyExists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	script, err := client.Script.Get(resourceData.Id())
	return script != nil, err
}
