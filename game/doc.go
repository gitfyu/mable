/*
Package game contains the logic for game features such as worlds and entities.

In order to run game logic in parallel without making development very complex, Mable uses a goroutine-per-world
model. This means that all logic within a single world, such as processing player actions, ticking entities,
or updating blocks are handled by a single goroutine, so that no synchronization is needed when the world or its
entities interact with other things in the same world. If you need to interact with a different world, you can
schedule a job to be executed in the other world by calling its World.Schedule function.

Functions that perform actions on a world or its contents should always be called from that worlds' handler
goroutine, unless a function's documentation explicitly states that it may be called concurrently.
*/
package game
