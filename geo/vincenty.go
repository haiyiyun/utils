package geo

import (
	"errors"
	"math"

	"github.com/haiyiyun/mongodb/geometry"
)

// these constants are used for vincentyDistance()
const a = 6378137
const b = 6356752.3142
const f = 1 / 298.257223563 // WGS-84 ellipsiod

/*
VincentyDistance computes the distances between two georgaphic coordinates
Args:
	p1: the 'starting' point, given in [0]longitude, [1]latitude as a PointCoordinates struct
	p2: the 'ending' point
Returns:
	A 2 element tuple: distance between the 2 points given in (1) meters
	The second element will return true upon a successful computation or
	false if the algorithm fails to converge. -1, false is returned upon failure
*/
func VincentyDistance(p1, p2 geometry.PointCoordinates) (float64, error) {
	// convert from degrees to radians
	var la1, lo1, la2, lo2 float64

	piRad := math.Pi / 180
	lo1 = p1[0] * piRad
	la1 = p1[1] * piRad
	lo2 = p2[0] * piRad
	la2 = p2[1] * piRad

	L := lo2 - lo1

	U1 := math.Atan((1 - f) * math.Tan(la1))
	U2 := math.Atan((1 - f) * math.Tan(la2))

	sinU1 := math.Sin(U1)
	cosU1 := math.Cos(U1)
	sinU2 := math.Sin(U2)
	cosU2 := math.Cos(U2)

	lambda := L
	lambdaP := 2 * math.Pi
	iterLimit := 20

	var sinLambda, cosLambda, sinSigma float64
	var cosSigma, sigma, sinAlpha, cosSqAlpha, cos2SigmaM, C float64

	for {
		if math.Abs(lambda-lambdaP) > 1e-12 && (iterLimit > 0) {
			iterLimit -= 1
		} else {
			break
		}
		sinLambda = math.Sin(lambda)
		cosLambda = math.Cos(lambda)

		sinSigma = math.Sqrt((cosU2*sinLambda)*(cosU2*sinLambda) + (cosU1*sinU2-sinU1*cosU2*cosLambda)*(cosU1*sinU2-sinU1*cosU2*cosLambda))
		if sinSigma == 0 {
			return 0, nil // co-incident points
		}

		cosSigma = sinU1*sinU2 + cosU1*cosU2*cosLambda
		sigma = math.Atan2(sinSigma, cosSigma)
		sinAlpha = cosU1 * cosU2 * sinLambda / sinSigma
		cosSqAlpha = 1 - sinAlpha*sinAlpha
		cos2SigmaM = cosSigma - 2*sinU1*sinU2/cosSqAlpha
		if math.IsNaN(cos2SigmaM) {
			cos2SigmaM = 0 // equatorial line: cosSqAlpha=0
		}

		C = f / 16 * cosSqAlpha * (4 + f*(4-3*cosSqAlpha))
		lambdaP = lambda
		lambda = L + (1-C)*f*sinAlpha*(sigma+C*sinSigma*(cos2SigmaM+C*cosSigma*(-1+2*cos2SigmaM*cos2SigmaM)))
	}
	if iterLimit == 0 {
		return -1, errors.New("vincenty algorithm failed to converge") // formula failed to converge
	}

	uSq := cosSqAlpha * (a*a - b*b) / (b * b)
	A := 1 + uSq/16384*(4096+uSq*(-768+uSq*(320-175*uSq)))
	B := uSq / 1024 * (256 + uSq*(-128+uSq*(74-47*uSq)))
	deltaSigma := B * sinSigma * (cos2SigmaM + B/4*(cosSigma*(-1+2*cos2SigmaM*cos2SigmaM)-B/6*cos2SigmaM*(-3+4*sinSigma*sinSigma)*(-3+4*cos2SigmaM*cos2SigmaM)))
	meters := b * A * (sigma - deltaSigma)
	return meters, nil
}
