package main 

fmt (
	"fmt"
)

type flag struct{
	sm int
	con int
}

var ma []int
var mf []flag
var n,m,sum int 

func di(i int ){
	if mf[i].sm != 0{
		ma[mf[i].sm]=ma[i]
	}
	if mf[i].con != 0{
		ma[mf[i].con] = !ma[i]
	}
}

func fin(i int){
	if ma[i] !=0 {
		di(i)
	}
	if(i==n){
		sum ++
		return
	}
	if 
	ma[i]=1
	di(i)
	fin(i+1)
	ma[i]=0 

}

func main(){
	var t,t1,t2 int
	fmt.Scanf("%d%d",&n,&m)
	ma = make([]int,n+1)
	mf = make([]flag,n+1)

	for i:=0;i<m;i++{
		fmt.Scanf("%d%d%d",&t,&t1,&t2)
		if t==2&&t1==t2{
			fmt.Println(0)
			return
		}
		if t1 > t2{
			var tx
			tx = t2
			t2 = t1
			t1 = tx
		}
		if t == 1 {
			mf[t1].sm = t2
		}else {
			mf[t1].con = t2
		}
	}
	sum = 0 
	fin(1)
}