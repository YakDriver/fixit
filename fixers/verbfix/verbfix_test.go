package verbfix

import (
	"fmt"
	"testing"
)

func TestSimpleReduce(t *testing.T) {
	before := fmt.Sprintf(`
func testAccAWSDBInstanceConfig_FinalSnapshotIdentifier(rInt int) string {
	return fmt.Sprintf(%[1]s
resource "aws_db_instance" "snapshot" {
  identifier = "tf-snapshot-%%d"
  final_snapshot_identifier = "foobarbaz-test-terraform-final-snapshot-%%d"
}
%[1]s, rInt, rInt)
}
`, "`")

	after := fmt.Sprintf(`
func testAccAWSDBInstanceConfig_FinalSnapshotIdentifier(rInt int) string {
	return fmt.Sprintf(%[1]s
resource "aws_db_instance" "snapshot" {
  identifier = "tf-snapshot-%%[1]d"
  final_snapshot_identifier = "foobarbaz-test-terraform-final-snapshot-%%[1]d"
}
%[1]s, rInt)
}
`, "`")

	if FixIt(before) != after {
		t.Errorf("got %v; want %v", FixIt(before), after)
	}
}

func TestDoubleReduce(t *testing.T) {
	before := fmt.Sprintf(`
func testAccAWSDBInstanceConfig_FinalSnapshotIdentifier(rInt int) string {
	return fmt.Sprintf(%[1]s
resource "aws_db_instance" "snapshot" {
  identifier = "tf-snapshot-%%d"
  final_snapshot_identifier = "%%s%%sfoobarbaz-test-terraform-final-snapshot-%%d"
}
%[1]s, rInt, "hello", "hello", rInt)
}
`, "`")

	after := fmt.Sprintf(`
func testAccAWSDBInstanceConfig_FinalSnapshotIdentifier(rInt int) string {
	return fmt.Sprintf(%[1]s
resource "aws_db_instance" "snapshot" {
  identifier = "tf-snapshot-%%[1]d"
  final_snapshot_identifier = "%%[2]s%%[2]sfoobarbaz-test-terraform-final-snapshot-%%[1]d"
}
%[1]s, rInt, "hello")
}
`, "`")

	if FixIt(before) != after {
		t.Errorf("got %v; want %v", FixIt(before), after)
	}
}

func TestMixedReduceable(t *testing.T) {
	before := fmt.Sprintf(`
func testAccAWSDBInstanceConfig_FinalSnapshotIdentifier(rInt int) string {
	return fmt.Sprintf(%[1]s
resource "aws_db_instance" "snapshot" {
  identifier = "tf-snapshot-%%[1]d"
  final_snapshot_identifier = "%%s%%sfoobarbaz-test-terraform-final-snapshot-%%[1]d"
  %%s
}
%[1]s, rInt, "hello", "hello", "bye")
}
`, "`")

	after := fmt.Sprintf(`
func testAccAWSDBInstanceConfig_FinalSnapshotIdentifier(rInt int) string {
	return fmt.Sprintf(%[1]s
resource "aws_db_instance" "snapshot" {
  identifier = "tf-snapshot-%%[1]d"
  final_snapshot_identifier = "%%[2]s%%[2]sfoobarbaz-test-terraform-final-snapshot-%%[1]d"
  %%[3]s
}
%[1]s, rInt, "hello", "bye")
}
`, "`")

	if FixIt(before) != after {
		t.Errorf("got %v; want %v", FixIt(before), after)
	}
}

func TestMultipleMatches(t *testing.T) {
	before := fmt.Sprintf(`
func testAccAWSDBInstanceConfig_FinalSnapshotIdentifier(rInt int) string {
	return fmt.Sprintf(%[1]s
resource "aws_db_instance" "snapshot" {
  identifier = "tf-snapshot-%%d"
  final_snapshot_identifier = "%%s%%sfoobarbaz-test-terraform-final-snapshot-%%d"
}
%[1]s, rInt, "hello", "hello", rInt)
}

func testAccAWSDBInstanceConfig_FinalSnapshotIdentifier2(rInt int) string {
	return fmt.Sprintf(%[1]s
resource "aws_db_instance" "snapshot" {
  identifier = "%%d-tf-%%d-snapshot-%%d"
  final_snapshot_identifier = "foobarbaz-test-terraform-final-snapshot-%%d"
}
%[1]s, rInt, rInt2, rInt, rInt2)
}
`, "`")

	after := fmt.Sprintf(`
func testAccAWSDBInstanceConfig_FinalSnapshotIdentifier(rInt int) string {
	return fmt.Sprintf(%[1]s
resource "aws_db_instance" "snapshot" {
  identifier = "tf-snapshot-%%[1]d"
  final_snapshot_identifier = "%%[2]s%%[2]sfoobarbaz-test-terraform-final-snapshot-%%[1]d"
}
%[1]s, rInt, "hello")
}

func testAccAWSDBInstanceConfig_FinalSnapshotIdentifier2(rInt int) string {
	return fmt.Sprintf(%[1]s
resource "aws_db_instance" "snapshot" {
  identifier = "%%[1]d-tf-%%[2]d-snapshot-%%[1]d"
  final_snapshot_identifier = "foobarbaz-test-terraform-final-snapshot-%%[2]d"
}
%[1]s, rInt, rInt2)
}
`, "`")

	if FixIt(before) != after {
		t.Errorf("got %v; want %v", FixIt(before), after)
	}
}

func TestMixedIndexNonIndex(t *testing.T) {
	before := fmt.Sprintf(`
func testAccAWSDBInstanceConfig_FinalSnapshotIdentifier(rInt int) string {
	return fmt.Sprintf(%[1]s
resource "aws_db_instance" "snapshot" {
  identifier = "tf-snapshot-%%[1]d"
  final_snapshot_identifier = "%%[1]d%%dfoobarbaz-test-terraform-final-snapshot"
}
%[1]s, rInt, rInt)
}
`, "`")

	after := fmt.Sprintf(`
func testAccAWSDBInstanceConfig_FinalSnapshotIdentifier(rInt int) string {
	return fmt.Sprintf(%[1]s
resource "aws_db_instance" "snapshot" {
  identifier = "tf-snapshot-%%[1]d"
  final_snapshot_identifier = "%%[1]d%%[1]dfoobarbaz-test-terraform-final-snapshot"
}
%[1]s, rInt)
}
`, "`")

	if FixIt(before) != after {
		t.Errorf("got %v; want %v", FixIt(before), after)
	}
}

func TestMixedIndexNonIndexFirst(t *testing.T) {
	before := fmt.Sprintf(`
func testAccAWSDBInstanceConfig_FinalSnapshotIdentifier(rInt int) string {
	return fmt.Sprintf(%[1]s
resource "aws_db_instance" "snapshot" {
  identifier = "tf-snapshot-%%d"
  final_snapshot_identifier = "%%[1]dfoobarbaz-test-terraform-final-snapshot"
}
%[1]s, rInt)
}
`, "`")

	after := fmt.Sprintf(`
func testAccAWSDBInstanceConfig_FinalSnapshotIdentifier(rInt int) string {
	return fmt.Sprintf(%[1]s
resource "aws_db_instance" "snapshot" {
  identifier = "tf-snapshot-%%[1]d"
  final_snapshot_identifier = "%%[1]dfoobarbaz-test-terraform-final-snapshot"
}
%[1]s, rInt)
}
`, "`")

	if FixIt(before) != after {
		t.Errorf("got %v; want %v", FixIt(before), after)
	}
}

func TestBackwardsIndexed(t *testing.T) {
	before := fmt.Sprintf(`
func testAccAWSDBInstanceConfig_FinalSnapshotIdentifier(rInt int) string {
	return fmt.Sprintf(%[1]s
resource "aws_db_instance" "snapshot" {
  identifier = "tf-snapshot-%%[3]d"
  final_snapshot_identifier = "%%[2]dfoobarbaz-%%[1]dtest-terraform-final-snapshot"
}
%[1]s, rInt1, rInt2, rInt3)
}
`, "`")

	after := fmt.Sprintf(`
func testAccAWSDBInstanceConfig_FinalSnapshotIdentifier(rInt int) string {
	return fmt.Sprintf(%[1]s
resource "aws_db_instance" "snapshot" {
  identifier = "tf-snapshot-%%d"
  final_snapshot_identifier = "%%dfoobarbaz-%%dtest-terraform-final-snapshot"
}
%[1]s, rInt3, rInt2, rInt1)
}
`, "`")

	if FixIt(before) != after {
		t.Errorf("got %v; want %v", FixIt(before), after)
	}
}

func TestFileFix(t *testing.T) {
	before, err := FileContent("tests/test_file_bad.txt")
	if err != nil {
		t.Error(err)
	}

	after, err := FileContent("tests/test_file_good.txt")
	if err != nil {
		t.Error(err)
	}

	if FixIt(before) != after {
		t.Errorf("got %v; want %v", FixIt(before), after)
	}
}

func TestTwoVerbsOneVar(t *testing.T) {
	before := fmt.Sprintf(`
func testAccAWSDBInstanceConfig_FinalSnapshotIdentifier(rInt int) string {
	return fmt.Sprintf(%[1]s
resource "aws_db_instance" "snapshot" {
  identifier = "tf-snapshot-%%[1]s"
  final_snapshot_identifier = %%[1]q
}
%[1]s, rName)
}
`, "`")

	after := fmt.Sprintf(`
func testAccAWSDBInstanceConfig_FinalSnapshotIdentifier(rInt int) string {
	return fmt.Sprintf(%[1]s
resource "aws_db_instance" "snapshot" {
  identifier = "tf-snapshot-%%[1]s"
  final_snapshot_identifier = %%[1]q
}
%[1]s, rName)
}
`, "`")

	if FixIt(before) != after {
		t.Errorf("got %v; want %v", FixIt(before), after)
	}
}
