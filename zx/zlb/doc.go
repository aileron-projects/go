// Package zlb provides load balancers.
//
// Provided algorithms:
//
//	|                    | Check  | Check  | Hash-base | Consistent | Computation |
//	| Algorithm          | weight | status | algorithm |    hash    | complexity  |
//	| ------------------ | ------ | ------ | --------- | ---------- | ----------- |
//	| Priority           |  Yes   |  Yes   |    No     |     --     |    O(n)     |
//	| Random             |  No    |  Yes   |    No     |     --     |    O(1)     |
//	| WeightedRandom     |  Yes   |  Yes   |    No     |     --     |    O(n)     |
//	| BasicRoundRobin    |  Yes   |  Yes   |    No     |     --     |    O(1)     |
//	| RoundRobin         |  Yes   |  Yes   |    No     |     --     |    O(n)     |
//	| RendezvousHash     |  Yes   |  Yes   |    Yes    |    Yes     |    O(n)     |
//	| JumpHash           |  No    |  Yes   |    Yes    |    Yes     |    O(1)     |
//	| DirectHash         |  Yes   |  Yes   |    Yes    |    No      |    O(1)     |
//	| WeightedDirectHash |  Yes   |  Yes   |    Yes    |    No      |    O(n)     |
//	| RingHash           |  Yes   |  Yes   |    Yes    |    Yes     |  O(log(n))  |
//	| Maglev             |  Yes   |  Yes   |    Yes    |    Yes     |    O(1)     |
//
// See the comments:
//   - [Priority]: priority based, or weight based load balancer.
//   - [Random]: random load balancer.
//   - [RandomW]: weighted random load balancer.
//   - [BasicRoundRobin]: basic round robin load balancer.
//   - [RoundRobin]: smooth round robin load balancer.
//   - [RendezvousHash]: rendezvous hash load balancer.
//   - [JumpHash]: jump hash load balancer.
//   - [DirectHash]: direct hash load balancer.
//   - [DirectHashW]: weighted direct hash load balancer.
//   - [RingHash]: ring hash load balancer.
//   - [Maglev]: maglev hash load balancer.
package zlb
