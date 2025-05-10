package collection

import (
	"testing"
)

func TestNew(t *testing.T) {
	coll := New[string, int]()
	if coll == nil {
		t.Fatal("Expected non-nil collection")
	}
	if len(coll.data) != 0 {
		t.Errorf("Expected empty data map, got %d elements", len(coll.data))
	}
}

func TestSet_Get_Has(t *testing.T) {
	coll := New[string, int]()
	key := "age"
	val := 25

	coll.Set(key, val)

	v, ok := coll.Get(key)
	if !ok {
		t.Error("Expected key to exist")
	}
	if v != val {
		t.Errorf("Expected value %v, got %v", val, v)
	}

	if !coll.Has(key) {
		t.Error("Has returned false for existing key")
	}

	if coll.Has("nonexistent") {
		t.Error("Has returned true for nonexistent key")
	}
}

func TestDelete(t *testing.T) {
	coll := New[string, int]()
	key := "age"
	val := 25
	coll.Set(key, val)

	coll.Delete(key)
	if _, ok := coll.Get(key); ok {
		t.Error("Key should be deleted")
	}

	coll.Delete("nonexistent") // не должно вызвать ошибок
}

func TestSize(t *testing.T) {
	coll := New[string, int]()
	if coll.Size() != 0 {
		t.Errorf("Expected size 0, got %d", coll.Size())
	}

	coll.Set("a", 1)
	coll.Set("b", 2)
	if coll.Size() != 2 {
		t.Errorf("Expected size 2, got %d", coll.Size())
	}
}

func TestKeys_Values(t *testing.T) {
	coll := New[string, int]()
	coll.Set("a", 1)
	coll.Set("b", 2)

	keys := coll.Keys()
	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(keys))
	}

	values := coll.Values()
	if len(values) != 2 {
		t.Errorf("Expected 2 values, got %d", len(values))
	}
}

func TestClear(t *testing.T) {
	coll := New[string, int]()
	coll.Set("a", 1)
	coll.Clear()
	if coll.Size() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", coll.Size())
	}
}

func TestEnsure(t *testing.T) {
	coll := New[string, int]()
	val := coll.Ensure("key", func(k string) int { return 42 })
	if val != 42 {
		t.Errorf("Expected default value 42, got %d", val)
	}

	val = coll.Ensure("key", func(k string) int { return 99 })
	if val != 42 {
		t.Errorf("Expected cached value 42, got %d", val)
	}
}

func TestHasAll_HasAny(t *testing.T) {
	coll := New[string, int]()
	coll.Set("a", 1)
	coll.Set("b", 2)

	if !coll.HasAll("a", "b") {
		t.Error("HasAll failed for all existing keys")
	}
	if coll.HasAll("a", "c") {
		t.Error("HasAll succeeded with missing key")
	}

	if !coll.HasAny("a", "c") {
		t.Error("HasAny failed when one key exists")
	}
	if coll.HasAny("x", "y") {
		t.Error("HasAny succeeded when no keys exist")
	}
}

func TestFirst_Last(t *testing.T) {
	coll := New[string, int]()
	coll.Set("a", 1)
	coll.Set("b", 2)
	coll.Set("c", 3)

	first := coll.First(1)
	if len(first) != 1 {
		t.Errorf("Expected 1 element, got %d", len(first))
	}

	last := coll.Last(1)
	if len(last) != 1 {
		t.Errorf("Expected 1 element, got %d", len(last))
	}
}

func TestAt_KeyAt(t *testing.T) {
	coll := New[string, int]()
	coll.Set("a", 1)
	coll.Set("b", 2)
	coll.Set("c", 3)

	val, err := coll.At(0)
	if err != nil || val != 1 && val != 2 && val != 3 {
		t.Errorf("At(0) failed: %v", err)
	}

	key, err := coll.KeyAt(0)
	if err != nil {
		t.Errorf("KeyAt(0) failed: %v", err)
	}
	if key != "a" && key != "b" && key != "c" {
		t.Errorf("Unexpected key: %v", key)
	}
}

func TestRandom(t *testing.T) {
	coll := New[string, int]()
	coll.Set("a", 1)
	coll.Set("b", 2)
	coll.Set("c", 3)

	random := coll.Random(2)
	if len(random) != 2 {
		t.Errorf("Expected 2 random elements, got %d", len(random))
	}
}

func TestFind_FindLast(t *testing.T) {
	coll := New[string, int]()
	coll.Set("a", 1)
	coll.Set("b", 2)
	coll.Set("c", 3)

	val, ok := coll.Find(func(v int, k string) bool { return v > 1 })
	if !ok || val <= 1 {
		t.Error("Find failed to find matching value")
	}

	val, ok = coll.FindLast(func(v int, k string) bool { return v < 3 })
	if !ok || val >= 3 {
		t.Error("FindLast failed to find matching value")
	}
}

func TestSweep(t *testing.T) {
	coll := New[string, int]()
	coll.Set("a", 1)
	coll.Set("b", 2)
	coll.Set("c", 3)

	count := coll.Sweep(func(v int, _ string) bool { return v < 3 })
	if count != 2 {
		t.Errorf("Expected 2 items swept, got %d", count)
	}
	if coll.Size() != 1 {
		t.Errorf("Expected 1 item left, got %d", coll.Size())
	}
}

func TestFilter_Partition(t *testing.T) {
	coll := New[string, int]()
	coll.Set("a", 1)
	coll.Set("b", 2)
	coll.Set("c", 3)

	filtered := coll.Filter(func(v int, _ string) bool { return v > 1 })
	if filtered.Size() != 2 {
		t.Errorf("Expected 2 filtered items, got %d", filtered.Size())
	}

	matched, unmatched := coll.Partition(func(v int, _ string) bool { return v > 1 })
	if matched.Size() != 2 || unmatched.Size() != 1 {
		t.Errorf("Partition mismatch: matched=%d, unmatched=%d", matched.Size(), unmatched.Size())
	}
}

func TestMap_Reduce(t *testing.T) {
	coll := New[string, int]()
	coll.Set("a", 1)
	coll.Set("b", 2)

	mapped := coll.Map(func(v int, _ string) int { return v * 2 })
	val, _ := mapped.Get("a")
	if val != 2 {
		t.Errorf("Mapped value should be 2, got %d", val)
	}

	sum := mapped.Reduce(func(acc, v int, _ string) int { return acc + v }, 0)
	if sum != 6 {
		t.Errorf("Reduced sum should be 6, got %d", sum)
	}
}

func TestSome_Every(t *testing.T) {
	coll := New[string, int]()
	coll.Set("a", 1)
	coll.Set("b", 2)

	if !coll.Some(func(v int, _ string) bool { return v == 1 }) {
		t.Error("Some failed to find match")
	}
	if coll.Every(func(v int, _ string) bool { return v == 1 }) {
		t.Error("Every incorrectly returned true")
	}
}

func TestEquals(t *testing.T) {
	coll1 := New[string, int]()
	coll1.Set("a", 1)
	coll1.Set("b", 2)

	coll2 := New[string, int]()
	coll2.Set("a", 1)
	coll2.Set("b", 2)

	if !coll1.Equals(coll2) {
		t.Error("Equals should return true for equal collections")
	}

	coll2.Set("b", 99)
	if coll1.Equals(coll2) {
		t.Error("Equals should return false for different collections")
	}
}

func TestSort_InPlace(t *testing.T) {
	coll := New[string, int]()
	coll.Set("a", 3)
	coll.Set("b", 1)
	coll.Set("c", 2)

	// Test Sort (should modify original collection)
	sorted := coll.Sort(func(v1, v2 int, _, _ string) bool { return v1 < v2 })

	// Verify sorted collection
	values := sorted.Values()
	if len(values) != 3 || values[0] != 1 || values[1] != 2 || values[2] != 3 {
		t.Errorf("Sort values not in order: %v", values)
	}

	// Verify it's the same collection (pointer equality)
	if sorted != coll {
		t.Errorf("Sort should return the same collection pointer")
	}
}

func TestSort_Empty(t *testing.T) {
	coll := New[string, int]()

	// Should not panic
	sorted := coll.Sort(func(v1, v2 int, _, _ string) bool { return v1 < v2 })

	if sorted.Size() != 0 {
		t.Errorf("Empty collection should remain empty after sort")
	}
}

func TestUnion_Intersection_Difference(t *testing.T) {
	coll1 := New[string, int]()
	coll1.Set("a", 1)
	coll1.Set("b", 2)

	coll2 := New[string, int]()
	coll2.Set("b", 20)
	coll2.Set("c", 3)

	union := coll1.Union(coll2)
	if union.Size() != 3 {
		t.Errorf("Union expected 3 elements, got %d", union.Size())
	}

	intersection := coll1.Intersection(coll2)
	if intersection.Size() != 1 {
		t.Errorf("Intersection expected 1 element, got %d", intersection.Size())
	}

	diff := coll1.Difference(coll2)
	if diff.Size() != 1 {
		t.Errorf("Difference expected 1 element, got %d", diff.Size())
	}
}

func TestSymmetricDifference(t *testing.T) {
	coll1 := New[string, int]()
	coll1.Set("a", 1)
	coll1.Set("b", 2)

	coll2 := New[string, int]()
	coll2.Set("b", 20)
	coll2.Set("c", 3)

	symDiff := coll1.SymmetricDifference(coll2)
	if symDiff.Size() != 2 {
		t.Errorf("Symmetric difference expected 2 elements, got %d", symDiff.Size())
	}
}

func TestClone(t *testing.T) {
	coll := New[string, int]()
	coll.Set("a", 1)
	coll.Set("b", 2)

	cloned := coll.Clone()
	if !coll.Equals(cloned) {
		t.Error("Cloned collection should be equal to original")
	}

	cloned.Set("c", 3)
	if coll.Equals(cloned) {
		t.Error("Original and modified clone should not be equal")
	}
}

func TestToReversed(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		coll := New[string, int]()
		coll.Set("a", 1)
		coll.Set("b", 2)
		coll.Set("c", 3)

		reversed := coll.ToReversed()

		// Verify size
		if reversed.Size() != 3 {
			t.Fatalf("Expected reversed size 3, got %d", reversed.Size())
		}

		// Get original values
		origValues := coll.Values()
		revValues := reversed.Values()

		// Verify order is reversed (may need sorting for consistent comparison)
		if len(origValues) > 0 && len(revValues) > 0 {
			if origValues[0] != revValues[len(revValues)-1] ||
				origValues[len(origValues)-1] != revValues[0] {
				t.Error("Values were not properly reversed")
			}
		}
	})

	t.Run("empty collection", func(t *testing.T) {
		empty := New[string, int]()
		reversed := empty.ToReversed()

		if reversed.Size() != 0 {
			t.Errorf("Expected empty reversed collection, got size %d", reversed.Size())
		}
	})
}
