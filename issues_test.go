package decimal

import (
	"math/big"
	"testing"
)

func TestIssue20(t *testing.T) {
	x := New(10240000000000, 0)
	x.mul(x, New(976563, 9))
	if v, _ := x.Int64(); v != 10000005120 {
		t.Fatal("error int64: ", v, x.Int(nil).Int64())
	}
}

func TestIssue65(t *testing.T) {
	const expected = "999999999000000000000000000000"
	r, _ := new(big.Rat).SetString(expected)
	r2 := new(Big).SetRat(r).Rat(nil)
	if r.Cmp(r2) != 0 {
		t.Fatalf("expected %q, got %q", r, r2)
	}
}

func TestIssue71(t *testing.T) {
	x, _ := new(Big).SetString("-433997231707950814777029946371807573425840064343095193931191306942897586882.200850175108941825587256711340679426793690849230895605323379098449524300541372392806145820741928")
	y := New(5, 0)
	ctx := Context{RoundingMode: ToZero, Precision: 364}

	z := new(Big)
	ctx.Quo(z, x, y)

	r, _ := new(Big).SetString("-86799446341590162955405989274361514685168012868619038786238261388579517376.4401700350217883651174513422681358853587381698461791210646758196899048601082744785612291641483856")
	if z.Cmp(r) != 0 || z.Scale() != r.Scale() {
		t.Fatalf(`Quo(%s, %s)
wanted: %s (%d)
got   : %s (%d)
`, x, y, r, -r.Scale(), z, -z.Scale())
	}
}

func TestIssue72(t *testing.T) {
	x, _ := new(Big).SetString("-8.45632792449080367076920780185655231664817924617196338687858969707575095137356626097186102468204972266270655439471710162223657196797091956190618036249568250856310899052975275153410779062120467574000771085625757386351708361318971283364474972153263288762761014798575650687906566E+474")
	y, _ := new(Big).SetString("4394389707820271499500265597691058417332780189928068605060129835915231607024733174128123086028964659911263805538968425927408117535552905751413991847682423230507052480632597367974353369255973450023914.06480266537851511912348920528179447782332532576762774258658035423323623047681531444628650113938865866058071268742035039370988065347285125745597527162817805470262344343643075954571122548882320506470664701832116848314413975179616459225485097673077072340532232446317251990415268245406080149594165531067657351225251495644780695372152557650401209918010537469259193951365404947434164664325966741900020673085975334136592934327584453217952431999450960191719318690339387778325911")
	z := new(Big)
	ctx := Context{Precision: 276}
	ctx.Rem(z, x, y)
	if !z.IsNaN(+1) {
		t.Fatalf(`Rem(%s, %s)
wanted: NaN
got   : %s
`, x, y, z)
	}
}
