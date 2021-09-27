package main

import (
	"math"
	"math/rand"
	"time"
)

type boid struct{
	position Vector2D
	velocity Vector2D
	id int
}

func (b *boid) start() {
	for {
		b.moveOne()
		time.Sleep(5 * time.Millisecond)
	}
}

func (b *boid) calcAcceleration() Vector2D {
	upper, lower := b.position.addValue(viewRadius), b.position.addValue(-viewRadius)
	avgVelocity := Vector2D{0, 0}
	count := 0.0
	
	for i := math.Max(lower.x, 0); i <= math.Min(upper.x, screenWidth); i++ {
		for j := math.Max(lower.y, 0); j <= math.Min(upper.y, screenHeight); j++ {
			if otherBoidId := boidMap[int(i)][int(j)]; otherBoidId != -1 && otherBoidId != b.id {
				if dist := boids[otherBoidId].position.distance(b.position); dist < viewRadius {
					count++
					avgVelocity = avgVelocity.add(boids[otherBoidId].velocity)
				}
			}
		}
	}
	accel := Vector2D{0, 0}

	if count > 0 {
		avgVelocity = avgVelocity.divideValue(count)
		accel = avgVelocity.subtract(b.velocity).multiplyValue(adjRate)
	}

	return accel
}

func (b *boid) moveOne() {
	b.velocity = b.velocity.add(b.calcAcceleration()).limit(-1, 1)
	boidMap[int(b.position.x)][int(b.position.y)] = -1
	b.position = b.position.add(b.velocity)
	boidMap[int(b.position.x)][int(b.position.y)] = b.id
	next := b.position.add(b.velocity)

	if next.x >= screenWidth || next.x < 0 {
		b.velocity = Vector2D{-b.velocity.x, b.velocity.y}
	}
	if next.y >= screenHeight || next.y < 0 {
		b.velocity = Vector2D{b.velocity.x, -b.velocity.y}
	}
}

func createBoid(bid int) {
	b := boid{
		position: Vector2D{rand.Float64() * screenWidth, rand.Float64() * screenHeight},
		velocity: Vector2D{(rand.Float64() * 2) - 1.0, (rand.Float64() * 2) - 1.0},
		id: bid,
	}

	boids[bid] = &b
	boidMap[int(b.position.x)][int(b.position.y)] = b.id
	go b.start()
}