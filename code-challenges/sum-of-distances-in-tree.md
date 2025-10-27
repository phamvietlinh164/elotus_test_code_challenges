# Sum of Distances in Tree

There is an undirected connected tree with n nodes labeled from 0 to n - 1 and n - 1 edges.

You are given the integer n and the array edges where edges[i] = [ai, bi] indicates that there is an edge between nodes ai and bi in the tree.

Return an array answer of length n where answer[i] is the sum of the distances between the ith node in the tree and all other nodes.

# Example 1:

![lc-sumdist](/images/lc-sumdist.jpg)

Input: n = 6, edges = [[0,1],[0,2],[2,3],[2,4],[2,5]]
Output: [8,12,6,10,10,10]
Explanation: The tree is shown above.
We can see that dist(0,1) + dist(0,2) + dist(0,3) + dist(0,4) + dist(0,5)
equals 1 + 1 + 2 + 2 + 2 = 8.
Hence, answer[0] = 8, and so on.

# Example 2:

Input: n = 1, edges = []
Output: [0]

# Example 3:

Input: n = 2, edges = [[1,0]]
Output: [1,1]

# Constraints:

- 1 <= n <= 3 * 104
- edges.length == n - 1
- edges[i].length == 2
- 0 <= ai, bi < n
- ai != bi
- The given input represents a valid tree.

# Go Template

```go
func sumOfDistancesInTree(n int, edges [][]int) []int {
    graph := make([][]int, n)
    for _, e := range edges {
        a, b := e[0], e[1]
        graph[a] = append(graph[a], b)
        graph[b] = append(graph[b], a)
    }

    res := make([]int, n)
    count := make([]int, n)
    for i := range count {
        count[i] = 1
    }

    // Step 1: Post-order DFS to compute res[0] and count[]
    var dfs1 func(int, int)
    dfs1 = func(node, parent int) {
        for _, nei := range graph[node] {
            if nei == parent {
                continue
            }
            dfs1(nei, node)
            count[node] += count[nei]
            res[node] += res[nei] + count[nei]
        }
    }

    // Step 2: Pre-order DFS to compute res[i] for all i
    var dfs2 func(int, int)
    dfs2 = func(node, parent int) {
        for _, nei := range graph[node] {
            if nei == parent {
                continue
            }
            res[nei] = res[node] - count[nei] + (n - count[nei])
            dfs2(nei, node)
        }
    }

    dfs1(0, -1)
    dfs2(0, -1)
    return res
}
```
