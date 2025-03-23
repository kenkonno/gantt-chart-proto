const AWS = require('aws-sdk');

exports.handler = async (event, context) => {
  console.log('Received event:', JSON.stringify(event, null, 2));

  const cloudfront = new AWS.CloudFront();
  console.log(event);
  const params = event;

  const invalidationParams = {
    DistributionId: params.distribution_id,
    InvalidationBatch: {
      CallerReference: new Date().getTime().toString(),
      Paths: {
        Quantity: params.paths.length,
        Items: params.paths
      }
    }
  };

  try {
    const response = await cloudfront.createInvalidation(invalidationParams).promise();
    console.log('Invalidation created successfully:', JSON.stringify(response, null, 2));
    return {
      statusCode: 200,
      body: JSON.stringify('Invalidation created successfully!')
    };
  } catch (error) {
    console.error('Error creating invalidation:', error);
    throw error;
  }
};
