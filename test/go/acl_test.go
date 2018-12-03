package _go

import (
	"encoding/xml"
	"github.com/journeymidnight/yig/api/datatype"
	. "github.com/journeymidnight/yig/test/go/lib"
	"net/http"
	"testing"
)

const (
	Case_PrivateAcl = `<AccessControlPolicy xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
		<Owner><ID>hehehehe</ID></Owner>
		<AccessControlList>
			<Grant>
				<Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="CanonicalUser">
					<ID>hehehehe</ID></Grantee>
				<Permission>FULL_CONTROL</Permission>
			</Grant>
			<Grant>
				<Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="Group">
					<URI>http://acs.amazonaws.com/groups/global/AllUsers</URI>
				</Grantee>
				<Permission>READ</Permission>
			</Grant>
		/AccessControlList>
	</AccessControlPolicy>`

	Case_PublicAcl = `<AccessControlPolicy xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
		<Owner>
			<ID>hehehehe</ID>
		</Owner>
		<AccessControlList>
			<Grant>
				<Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="CanonicalUser">
					<ID>hehehehe</ID>
				</Grantee>
				<Permission>FULL_CONTROL</Permission>
			</Grant>
			<Grant>
				<Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="Group">
					<URI>http://acs.amazonaws.com/groups/global/AllUsers</URI>
				</Grantee>
				<Permission>READ</Permission>
			</Grant>
		</AccessControlList>
	</AccessControlPolicy>`
)

func Test_ACL_Prepare(t *testing.T) {
	sc := NewS3()
	err := sc.MakeBucket(TEST_BUCKET)
	if err != nil {
		t.Fatal("MakeBucket err:", err)
		panic(err)
	}
	t.Log("MakeBucket Success.")
	err = sc.PutObject(TEST_BUCKET, TEST_KEY, TEST_VALUE)
	if err != nil {
		t.Fatal("PutObject err:", err)
		panic(err)
	}
	t.Log("PutObject Success.")
}

func Test_PutBucketAcl(t *testing.T) {
	sc := NewS3()
	err := sc.PutBucketAcl(TEST_BUCKET, BucketCannedACLPrivate)
	if err != nil {
		t.Fatal("PutBucketAcl err:", err)
	}
	t.Log("PutBucketAcl Success!")
}

func Test_GetBucketAcl(t *testing.T) {
	sc := NewS3()
	out, err := sc.GetBucketAcl(TEST_BUCKET)
	if err != nil {
		t.Fatal("GetBucketAcl err:", err)
	}
	t.Log("GetBucketAcl Success! out:", out)
}

func Test_PutObjectAclWithNoSuchObject(t *testing.T) {
	sc := NewS3()
	err := sc.PutObjectAcl(TEST_BUCKET, TEST_KEY+"NONE", BucketCannedACLPrivate)
	if err == nil {
		t.Fatal("Test_PutObjectAclWithNoSuchObject err: We have no such key", TEST_KEY)
	}
	t.Log("PutObjectAclWithNoSuchObject Success!")
}

func Test_PutObjectAcl(t *testing.T) {
	sc := NewS3()
	err := sc.PutObjectAcl(TEST_BUCKET, TEST_KEY, BucketCannedACLPrivate)
	if err != nil {
		t.Fatal("PutBucketAcl err:", err)
	}
	t.Log("PutBucketAcl Success!")
}

func Test_GetObjectAcl(t *testing.T) {
	sc := NewS3()
	out, err := sc.GetObjectAcl(TEST_BUCKET, TEST_KEY)
	if err != nil {
		t.Fatal("GetObjectAcl err:", err)
	}
	t.Log("GetObjectAcl Success! out:", out)
}

func Test_PutObjectPublicAclWithXml(t *testing.T) {
	sc := NewS3()
	url := GenTestObjectUrl(sc)

	statusCode, _, err := HTTPRequestToGetObject(url)
	if err != nil {
		t.Fatal("GetObject err:", err)
	}
	//StatusCode should be AccessDenied
	if statusCode != http.StatusForbidden {
		t.Fatal("StatusCode should be AccessDenied(403), but the code is:", statusCode)
	}
	t.Log("GetObject Without ACL test Success.")

	//set policy
	var policy = &datatype.AccessControlPolicy{}
	err = xml.Unmarshal([]byte(Case_PublicAcl), policy)
	if err != nil {
		t.Fatal("PutObjectPublicAclWithXml err:", err)
	}
	acl := TransferToS3AccessControlPolicy(policy)
	if acl == nil {
		t.Fatal("PutObjectPublicAclWithXml err:", "empty acl!")
	}
	err = sc.PutObjectAclWithXml(TEST_BUCKET, TEST_KEY, acl)
	if err != nil {
		t.Fatal("PutObjectAclWithXml err:", err)
	}
	t.Log("PutObjectAclWithXml Success!")

	out, err := sc.GetObjectAcl(TEST_BUCKET, TEST_KEY)
	if err != nil {
		t.Fatal("GetObjectAcl err:", err)
	}
	t.Log("GetObjectAcl Success! out:", out)

	// After set acl
	statusCode, data, err := HTTPRequestToGetObject(url)
	if err != nil {
		t.Fatal("GetObject err:", err)
	}
	//StatusCode should be STATUS_OK
	if statusCode != http.StatusOK {
		t.Fatal("StatusCode should be STATUS_OK(200), but the code is:", statusCode)
	}
	t.Log("Get object value:", string(data))

}

func Test_ACL_End(t *testing.T) {
	sc := NewS3()
	err := sc.DeleteObject(TEST_BUCKET, TEST_KEY)
	if err != nil {
		t.Log("DeleteObject err:", err)
	}
	err = sc.DeleteBucket(TEST_BUCKET)
	if err != nil {
		t.Fatal("DeleteBucket err:", err)
		panic(err)
	}
}
