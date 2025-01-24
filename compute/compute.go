package compute

import (
	"context"
	"fmt"
	"math/big"
)

type splitResult struct {
	p, q, r *big.Int
}

func newSplitResult() *splitResult {
	return &splitResult{
		p: new(big.Int),
		q: new(big.Int),
		r: new(big.Int),
	}
}

func binarySplit(ctx context.Context, a, b int64) *splitResult {
	select {
	case <-ctx.Done():
		fmt.Println("Context cancelled inside binarySplit")
		return nil
	default:
	}

	result := newSplitResult()

	if b == a+1 {
		// P = -(6a-5)(2a-1)(6a-1)
		p1 := big.NewInt(6*a - 5)
		p2 := big.NewInt(2*a - 1)
		p3 := big.NewInt(6*a - 1)
		result.p.Mul(p1, p2)
		result.p.Mul(result.p, p3)
		result.p.Neg(result.p)

		// Q = 10939058860032000 * a^3
		base := big.NewInt(10939058860032000)
		exp := big.NewInt(a)
		result.q.Exp(exp, big.NewInt(3), nil)
		result.q.Mul(result.q, base)

		// R = P * (545140134*a + 13591409)
		term := big.NewInt(545140134)
		term.Mul(term, big.NewInt(a))
		term.Add(term, big.NewInt(13591409))
		result.r.Mul(result.p, term)
	} else {
		m := (a + b) / 2
		left := binarySplit(ctx, a, m)
		if left == nil {
			return nil
		}
		right := binarySplit(ctx, m, b)
		if right == nil {
			return nil
		}

		// P = Pam * Pmb
		result.p.Mul(left.p, right.p)

		// Q = Qam * Qmb
		result.q.Mul(left.q, right.q)

		// R = Qmb * Ram + Pam * Rmb
		term1 := new(big.Int).Mul(right.q, left.r)
		term2 := new(big.Int).Mul(left.p, right.r)
		result.r.Add(term1, term2)
	}

	return result
}

// ComputePi calculates Pi to the specified number of iterations
func ComputePi(ctx context.Context, iterations int64) *big.Float {
	select {
	case <-ctx.Done():
		fmt.Println("Context cancelled inside ComputePi")
		return nil
	default:
	}
	split := binarySplit(ctx, 1, iterations+1)
	if split == nil {
		return nil
	}

	// Calculate C = 426880 * sqrt(10005)
	c := big.NewFloat(426880)
	c.Mul(c, sqrt(big.NewFloat(10005)))

	// Calculate denominator = 13591409*Q + R
	denom := new(big.Int)
	term := new(big.Int).Mul(big.NewInt(13591409), split.q)
	denom.Add(term, split.r)

	// Final result = (C * Q) / denominator
	result := new(big.Float).SetInt(split.q)
	result.Mul(result, c)
	denomFloat := new(big.Float).SetInt(denom)
	result.Quo(result, denomFloat)

	return result
}

func sqrt(x *big.Float) *big.Float {
	// Newton's method for square root
	z := new(big.Float).Set(x)
	t := new(big.Float)
	half := big.NewFloat(0.5)

	for i := 0; i < 10; i++ {
		t.Quo(x, z)
		t.Add(t, z)
		t.Mul(t, half)
		if z.Cmp(t) == 0 {
			break
		}
		z.Set(t)
	}
	return z
}
