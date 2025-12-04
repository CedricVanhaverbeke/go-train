package route

// https://www.youtube.com/watch?time_continue=1&v=VcCLGNdSusk&embeds_referring_euri=https%3A%2F%2Fchatgpt.com%2F&source_ve_path=MjM4NTE
var example = `
	<?xml version="1.0" encoding="UTF-8"?>
<gpx creator="StravaGPX" version="1.1" xmlns="http://www.topografix.com/GPX/1/1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd">
 <metadata>
  <name>cedric</name>
  <author>
   <name>Cedric Vanhaverbeke</name>
   <link href="https://www.strava.com/athletes/39174470"/>
  </author>
  <copyright author="OpenStreetMap contributors">
   <year>2020</year>
   <license>https://www.openstreetmap.org/copyright</license>
  </copyright>
  <link href="https://www.strava.com/routes/3312769391928190254"/>
 </metadata>
 <trk>
  <name>cedric</name>
  <link href="https://www.strava.com/routes/3312769391928190254"/>
  <type>walking</type>
  <trkseg>
   <trkpt lat="38.90724" lon="-77.02964">
    <ele>27.770000000000003</ele>
   </trkpt>
   <trkpt lat="38.907240005101336" lon="-77.03041333327305">
    <ele>27.590000000000003</ele>
   </trkpt>
   <trkpt lat="38.90724000510134" lon="-77.03118666654615">
    <ele>27.580000000000002</ele>
   </trkpt>
   <trkpt lat="38.90724" lon="-77.03196000000001">
    <ele>28.110000000000003</ele>
   </trkpt>
   <trkpt lat="38.906650000000006" lon="-77.03196000000001">
    <ele>27.39</ele>
   </trkpt>
   <trkpt lat="38.90641" lon="-77.03192">
    <ele>26.990000000000002</ele>
   </trkpt>
   <trkpt lat="38.905820000000006" lon="-77.03157">
    <ele>26.19</ele>
   </trkpt>
   <trkpt lat="38.90567" lon="-77.03151000000001">
    <ele>26.01</ele>
   </trkpt>
   <trkpt lat="38.90567" lon="-77.0313">
    <ele>25.880000000000003</ele>
   </trkpt>
   <trkpt lat="38.9056" lon="-77.03129000000001">
    <ele>25.78</ele>
   </trkpt>
   <trkpt lat="38.90559000293838" lon="-77.03045999990319">
    <ele>26.09</ele>
   </trkpt>
   <trkpt lat="38.90558" lon="-77.02963000000001">
    <ele>26.650000000000002</ele>
   </trkpt>
   <trkpt lat="38.90557500249633" lon="-77.02886500000109">
    <ele>26.88</ele>
   </trkpt>
   <trkpt lat="38.905570000000004" lon="-77.02810000000001">
    <ele>27.770000000000003</ele>
   </trkpt>
   <trkpt lat="38.90596" lon="-77.02810000000001">
    <ele>28.05</ele>
   </trkpt>
   <trkpt lat="38.90596" lon="-77.02817">
    <ele>27.87</ele>
   </trkpt>
   <trkpt lat="38.906470000000006" lon="-77.02818">
    <ele>28.3</ele>
   </trkpt>
   <trkpt lat="38.906470000000006" lon="-77.02901">
    <ele>27.020000000000003</ele>
   </trkpt>
   <trkpt lat="38.906470000000006" lon="-77.02902">
    <ele>27.020000000000003</ele>
   </trkpt>
   <trkpt lat="38.906470000000006" lon="-77.02946">
    <ele>27.240000000000002</ele>
   </trkpt>
   <trkpt lat="38.90571500003168" lon="-77.02947000010593">
    <ele>26.64</ele>
   </trkpt>
   <trkpt lat="38.90496" lon="-77.02948">
    <ele>27.1</ele>
   </trkpt>
   <trkpt lat="38.904790000000006" lon="-77.02945000000001">
    <ele>27.080000000000002</ele>
   </trkpt>
   <trkpt lat="38.90466" lon="-77.02908000000001">
    <ele>27.44</ele>
   </trkpt>
   <trkpt lat="38.904360000000004" lon="-77.02822">
    <ele>28.05</ele>
   </trkpt>
   <trkpt lat="38.90444" lon="-77.02822">
    <ele>28.14</ele>
   </trkpt>
   <trkpt lat="38.90447" lon="-77.02820000000001">
    <ele>28.14</ele>
   </trkpt>
   <trkpt lat="38.904430000000005" lon="-77.02810000000001">
    <ele>28.12</ele>
   </trkpt>
   <trkpt lat="38.90431" lon="-77.02810000000001">
    <ele>28.05</ele>
   </trkpt>
   <trkpt lat="38.90415" lon="-77.02765000000001">
    <ele>27.090000000000003</ele>
   </trkpt>
   <trkpt lat="38.904160000000005" lon="-77.02737">
    <ele>26.6</ele>
   </trkpt>
   <trkpt lat="38.904160000000005" lon="-77.02738000000001">
    <ele>26.6</ele>
   </trkpt>
   <trkpt lat="38.904160000000005" lon="-77.02720000000001">
    <ele>26.27</ele>
   </trkpt>
   <trkpt lat="38.90401000000001" lon="-77.02720000000001">
    <ele>26.1</ele>
   </trkpt>
   <trkpt lat="38.90380000149101" lon="-77.02659499800868">
    <ele>24.82</ele>
   </trkpt>
   <trkpt lat="38.90359" lon="-77.02599000000001">
    <ele>23.46</ele>
   </trkpt>
   <trkpt lat="38.90346" lon="-77.02599000000001">
    <ele>23.1</ele>
   </trkpt>
   <trkpt lat="38.90421599999303" lon="-77.02598800008522">
    <ele>24.540000000000003</ele>
   </trkpt>
   <trkpt lat="38.904971999986" lon="-77.02598600012783">
    <ele>26.160000000000004</ele>
   </trkpt>
   <trkpt lat="38.90572799997895" lon="-77.02598400012785">
    <ele>27.26</ele>
   </trkpt>
   <trkpt lat="38.90648399997186" lon="-77.02598200008529">
    <ele>28.330000000000002</ele>
   </trkpt>
   <trkpt lat="38.90724" lon="-77.02598">
    <ele>29.380000000000003</ele>
   </trkpt>
   <trkpt lat="38.90648000009168" lon="-77.02598">
    <ele>28.330000000000002</ele>
   </trkpt>
   <trkpt lat="38.90572" lon="-77.02598">
    <ele>27.26</ele>
   </trkpt>
   <trkpt lat="38.90572" lon="-77.02607">
    <ele>27.39</ele>
   </trkpt>
   <trkpt lat="38.90565" lon="-77.02608000000001">
    <ele>27.31</ele>
   </trkpt>
   <trkpt lat="38.90565" lon="-77.02689000000001">
    <ele>28.130000000000003</ele>
   </trkpt>
   <trkpt lat="38.905150000038184" lon="-77.02690000006967">
    <ele>27.560000000000002</ele>
   </trkpt>
   <trkpt lat="38.904650000000004" lon="-77.02691">
    <ele>26.910000000000004</ele>
   </trkpt>
   <trkpt lat="38.904650000000004" lon="-77.02704">
    <ele>26.96</ele>
   </trkpt>
   <trkpt lat="38.90395" lon="-77.02704">
    <ele>25.86</ele>
   </trkpt>
   <trkpt lat="38.903616674588505" lon="-77.02607665763868">
    <ele>23.71</ele>
   </trkpt>
   <trkpt lat="38.9032833412612" lon="-77.02511332432293">
    <ele>22.27</ele>
   </trkpt>
   <trkpt lat="38.902950000000004" lon="-77.02415">
    <ele>21.23</ele>
   </trkpt>
   <trkpt lat="38.90306" lon="-77.02416000000001">
    <ele>21.37</ele>
   </trkpt>
   <trkpt lat="38.9031" lon="-77.02413">
    <ele>21.37</ele>
   </trkpt>
   <trkpt lat="38.903180000000006" lon="-77.0241">
    <ele>21.37</ele>
   </trkpt>
   <trkpt lat="38.904070000000004" lon="-77.02411000000001">
    <ele>24.040000000000003</ele>
   </trkpt>
   <trkpt lat="38.904050000000005" lon="-77.02398000000001">
    <ele>23.880000000000003</ele>
   </trkpt>
   <trkpt lat="38.90484749999577" lon="-77.02398000000001">
    <ele>24.92</ele>
   </trkpt>
   <trkpt lat="38.90564499999153" lon="-77.02398000000001">
    <ele>25.84</ele>
   </trkpt>
   <trkpt lat="38.90644249998729" lon="-77.02398000000001">
    <ele>26.7</ele>
   </trkpt>
   <trkpt lat="38.90724" lon="-77.02398000000001">
    <ele>27.69</ele>
   </trkpt>
   <trkpt lat="38.90724" lon="-77.02389000000001">
    <ele>27.680000000000003</ele>
   </trkpt>
   <trkpt lat="38.90723500413869" lon="-77.02290500006131">
    <ele>27.51</ele>
   </trkpt>
   <trkpt lat="38.907230000000006" lon="-77.02192000000001">
    <ele>27.96</ele>
   </trkpt>
   <trkpt lat="38.90647999982031" lon="-77.02192000000001">
    <ele>27.39</ele>
   </trkpt>
   <trkpt lat="38.905730000000005" lon="-77.02192000000001">
    <ele>27.13</ele>
   </trkpt>
   <trkpt lat="38.90565" lon="-77.02192000000001">
    <ele>27.13</ele>
   </trkpt>
   <trkpt lat="38.90564500322813" lon="-77.02279000006583">
    <ele>26.19</ele>
   </trkpt>
   <trkpt lat="38.905640000000005" lon="-77.02366">
    <ele>25.86</ele>
   </trkpt>
   <trkpt lat="38.905640002763796" lon="-77.02285499994494">
    <ele>26.11</ele>
   </trkpt>
   <trkpt lat="38.905640000000005" lon="-77.02205000000001">
    <ele>26.94</ele>
   </trkpt>
   <trkpt lat="38.9049699999518" lon="-77.02205500004756">
    <ele>25.150000000000002</ele>
   </trkpt>
   <trkpt lat="38.904300000000006" lon="-77.02206000000001">
    <ele>23.13</ele>
   </trkpt>
   <trkpt lat="38.90429" lon="-77.02192000000001">
    <ele>23.12</ele>
   </trkpt>
   <trkpt lat="38.90362000007713" lon="-77.02192000000001">
    <ele>21.84</ele>
   </trkpt>
   <trkpt lat="38.902950000000004" lon="-77.02192000000001">
    <ele>20.19</ele>
   </trkpt>
   <trkpt lat="38.90319333772517" lon="-77.02120333821381">
    <ele>19.85</ele>
   </trkpt>
   <trkpt lat="38.903436671069386" lon="-77.02048667151516">
    <ele>19.24</ele>
   </trkpt>
   <trkpt lat="38.90368" lon="-77.01977000000001">
    <ele>18.770000000000003</ele>
   </trkpt>
   <trkpt lat="38.90377" lon="-77.01978000000001">
    <ele>18.97</ele>
   </trkpt>
   <trkpt lat="38.90379" lon="-77.01979">
    <ele>18.97</ele>
   </trkpt>
   <trkpt lat="38.90411" lon="-77.01977000000001">
    <ele>20.1</ele>
   </trkpt>
   <trkpt lat="38.90484499990384" lon="-77.01977999989518">
    <ele>23.57</ele>
   </trkpt>
   <trkpt lat="38.90558" lon="-77.01979">
    <ele>25.61</ele>
   </trkpt>
   <trkpt lat="38.90558" lon="-77.01976">
    <ele>25.580000000000002</ele>
   </trkpt>
   <trkpt lat="38.905640000000005" lon="-77.01976">
    <ele>25.77</ele>
   </trkpt>
   <trkpt lat="38.90558" lon="-77.01976">
    <ele>25.580000000000002</ele>
   </trkpt>
   <trkpt lat="38.90558" lon="-77.01979">
    <ele>25.61</ele>
   </trkpt>
   <trkpt lat="38.90484500009702" lon="-77.01977999989779">
    <ele>23.57</ele>
   </trkpt>
   <trkpt lat="38.90411" lon="-77.01977000000001">
    <ele>20.1</ele>
   </trkpt>
   <trkpt lat="38.90379" lon="-77.01979">
    <ele>18.97</ele>
   </trkpt>
   <trkpt lat="38.90377" lon="-77.01978000000001">
    <ele>18.97</ele>
   </trkpt>
   <trkpt lat="38.90368" lon="-77.01977000000001">
    <ele>18.770000000000003</ele>
   </trkpt>
   <trkpt lat="38.903690000000005" lon="-77.01962">
    <ele>18.7</ele>
   </trkpt>
   <trkpt lat="38.90357" lon="-77.01957">
    <ele>18.410000000000004</ele>
   </trkpt>
   <trkpt lat="38.90373" lon="-77.01908">
    <ele>18.18</ele>
   </trkpt>
   <trkpt lat="38.903710000000004" lon="-77.01905000000001">
    <ele>18.13</ele>
   </trkpt>
   <trkpt lat="38.90377" lon="-77.01893000000001">
    <ele>18.060000000000002</ele>
   </trkpt>
   <trkpt lat="38.90384" lon="-77.01880000000001">
    <ele>18.0</ele>
   </trkpt>
   <trkpt lat="38.90399" lon="-77.01838000000001">
    <ele>18.230000000000004</ele>
   </trkpt>
   <trkpt lat="38.904030000000006" lon="-77.01828">
    <ele>18.330000000000002</ele>
   </trkpt>
   <trkpt lat="38.90401000000001" lon="-77.01827">
    <ele>18.330000000000002</ele>
   </trkpt>
   <trkpt lat="38.90435000418376" lon="-77.0172800047303">
    <ele>19.270000000000003</ele>
   </trkpt>
   <trkpt lat="38.90469" lon="-77.01629000000001">
    <ele>20.130000000000003</ele>
   </trkpt>
   <trkpt lat="38.90444" lon="-77.01704000000001">
    <ele>19.69</ele>
   </trkpt>
   <trkpt lat="38.90449" lon="-77.01707">
    <ele>19.720000000000002</ele>
   </trkpt>
   <trkpt lat="38.90424" lon="-77.01782">
    <ele>18.72</ele>
   </trkpt>
   <trkpt lat="38.90408" lon="-77.01836">
    <ele>18.53</ele>
   </trkpt>
   <trkpt lat="38.90408" lon="-77.01844000000001">
    <ele>18.53</ele>
   </trkpt>
   <trkpt lat="38.903960000000005" lon="-77.01882">
    <ele>18.340000000000003</ele>
   </trkpt>
   <trkpt lat="38.90417" lon="-77.01882">
    <ele>18.99</ele>
   </trkpt>
   <trkpt lat="38.904250000000005" lon="-77.01882">
    <ele>19.32</ele>
   </trkpt>
   <trkpt lat="38.90433" lon="-77.01886">
    <ele>19.75</ele>
   </trkpt>
   <trkpt lat="38.90495499992296" lon="-77.01886">
    <ele>21.96</ele>
   </trkpt>
   <trkpt lat="38.90558" lon="-77.01886">
    <ele>24.310000000000002</ele>
   </trkpt>
   <trkpt lat="38.90558" lon="-77.01883000000001">
    <ele>24.450000000000003</ele>
   </trkpt>
   <trkpt lat="38.90565" lon="-77.01883000000001">
    <ele>24.540000000000003</ele>
   </trkpt>
   <trkpt lat="38.905710000000006" lon="-77.01883000000001">
    <ele>24.55</ele>
   </trkpt>
   <trkpt lat="38.90572" lon="-77.01885">
    <ele>24.57</ele>
   </trkpt>
   <trkpt lat="38.906400000000005" lon="-77.01885">
    <ele>25.05</ele>
   </trkpt>
   <trkpt lat="38.90644" lon="-77.01885">
    <ele>25.120000000000005</ele>
   </trkpt>
   <trkpt lat="38.90644334008986" lon="-77.01796000010553">
    <ele>25.84</ele>
   </trkpt>
   <trkpt lat="38.9064466734231" lon="-77.0170700001274">
    <ele>26.070000000000004</ele>
   </trkpt>
   <trkpt lat="38.90645000000001" lon="-77.01618">
    <ele>25.87</ele>
   </trkpt>
  </trkseg>
 </trk>
</gpx>

`
