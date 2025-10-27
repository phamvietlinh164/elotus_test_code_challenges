# Maximum Length of Repeated Subarray

Given two integer arrays **nums1** and **nums2**, return the maximum length of a subarray that appears in **both** arrays.

# Example 1:

Input: nums1 = [1,2,3,2,1], nums2 = [3,2,1,4,7]
Output: 3
Explanation: The repeated subarray with maximum length is [3,2,1].

# Example 2:

Input: nums1 = [0,0,0,0,0], nums2 = [0,0,0,0,0]
Output: 5

# Constraints:

  - **1 <= nums1.length, nums2.length <= 1000**

  - **0 <= nums1[i], nums2[i] <= 100**

# Go Template

```go
func findLength(nums1 []int, nums2 []int) int {
    m, n := len(nums1), len(nums2)
    prev := make([]int, n+1)
    curr := make([]int, n+1)
    maxLen := 0

    for i := 1; i <= m; i++ {
        for j := 1; j <= n; j++ {
            if nums1[i-1] == nums2[j-1] {
                curr[j] = prev[j-1] + 1
                if curr[j] > maxLen {
                    maxLen = curr[j]
                }
            } else {
                curr[j] = 0
            }
        }
        prev, curr = curr, make([]int, n+1)
    }
    return maxLen
}
```
